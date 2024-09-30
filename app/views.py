from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
import pandas as pd
import numpy as np
import json
import os

# Path to the folder where your files are stored
BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

@csrf_exempt
def fake_deletion(request):
    if request.method == "POST":
        try:
            data = json.loads(request.body)
            # Convert relative paths to absolute paths
            input_path = os.path.join(BASE_DIR, data.get("input_path"))
            output_path = os.path.join(BASE_DIR, data.get("output_path"))

            if input_path and output_path:
                try:
                    # Read the CSV file from the input path
                    real_data = pd.read_csv(input_path)
                    # Process the data
                    fake_processed_data = introduce_missing_values(real_data, missing_rate=0.1)
                    # Save the processed data to the output path
                    fake_processed_data.to_csv(output_path, index=False)
                    return JsonResponse({"status": "success", "message": f"File saved at {output_path}"})
                except Exception as e:
                    return JsonResponse({"status": "error", "message": str(e)})
            else:
                return JsonResponse({"status": "error", "message": "Invalid input/output paths"})
        except Exception as e:
            return JsonResponse({"status": "error", "message": str(e)})
    return JsonResponse({"status": "error", "message": "Invalid request method"})

def introduce_missing_values(data, missing_rate=0.1):
    np.random.seed(42)
    data = data.copy()
    mask = np.random.rand(*data.shape) < missing_rate
    data[mask] = np.nan
    return data
