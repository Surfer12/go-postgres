from flask import Flask, request, jsonify
from google.cloud import aiplatform
from jsonschema import validate, ValidationError
import os
import logging

app = Flask(__name__)
logging.basicConfig(level=logging.INFO)

# Configuration
PROJECT_ID = os.environ.get("PROJECT_ID")  # GCP project ID
REGION = os.environ.get("REGION")          # GCP region (e.g., "us-central1")
ENDPOINT_ID = os.environ.get("ENDPOINT_ID") # Vertex AI Endpoint ID

# Initialize Vertex AI client
aiplatform.init(project=PROJECT_ID, location=REGION)
endpoint_client = aiplatform.Endpoint(endpoint_name=ENDPOINT_ID)

# Define schema for input validation
schema = {
    "type": "object",
    "properties": {
        "instances": {
            "type": "array",
            "items": {
                "type": "array",
                "items": {"type": "number"}
            }
        }
    },
    "required": ["instances"]
}

@app.route('/predict', methods=['POST'])
def predict():
    try:
        request_json = request.get_json()
        if not request_json:
            return jsonify({"error": "Request must be in JSON format", "code": "INVALID_FORMAT"}), 400

        try:
            validate(instance=request_json, schema=schema)
        except ValidationError as ve:
            return jsonify({"error": "Invalid input format", "code": "INVALID_INPUT", "details": str(ve)}), 400

        instances = request_json["instances"]
        try:
            prediction = endpoint_client.predict(instances=instances)
            return jsonify(prediction.predictions)
        except GoogleAPIError as e:
            logging.error(f"Vertex AI error: {e}")
            return jsonify({
                "error": "Vertex AI prediction failed",
                "code": "VERTEX_AI_ERROR",
                "details": str(e)
            }), 500
    except Exception as e:
        logging.error(f"Unexpected error: {e}")
        return jsonify({
            "error": "Internal server error",
            "code": "INTERNAL_ERROR",
            "details": "An unexpected error occurred. Please try again later."
        }), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=int(os.environ.get("PORT", 8080))) 