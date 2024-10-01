from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
import pandas as pd
import numpy as np
import json
import os
from sklearn.impute import KNNImputer
from sklearn.linear_model import LinearRegression
from sklearn.experimental import enable_iterative_imputer
from sklearn.impute import IterativeImputer
from sklearn.preprocessing import StandardScaler
from scipy.interpolate import PchipInterpolator
from sklearn.decomposition import TruncatedSVD

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







@csrf_exempt
def data_imputation(request):
    if request.method == "POST":
        print("Request Body (raw):", request.body)  # Ham veri olarak body'yi yazdırır
        print("Request Method:", request.method)  # HTTP metodunu yazdırır
        print("Request Headers:", request.headers)  # Başlıkları yazdırır
        try:
            
            data = json.loads(request.body)
            # Convert relative paths to absolute paths
            input_path = os.path.join(BASE_DIR, data.get("input_path"))
            output_path = os.path.join(BASE_DIR, data.get("output_path"))
            kwargs = data.get("kwargs", {})
            answer = data.get("answer", None)  # 'answer' key'ini alıyoruz
            if input_path and output_path:
                try:
                    # Read the CSV file from the input path
                    real_data = pd.read_csv(input_path)

                    # Choose imputation method based on the answer
                    if answer == "constant":
                        print("c")

                        result_data = fill_with_constant(real_data, **kwargs)
                    elif answer == "mean":
                        print("d")
                        result_data = fill_with_mean_or_median(real_data, method="mean")
                    elif answer == "median":
                        print("e")
                        result_data = fill_with_mean_or_median(real_data, method="median")
                    elif answer == "knn":
                        print("f")
                        result_data = fill_with_knn(real_data, **kwargs)
                    elif answer == "linear_regression":
                        print("g")
                        result_data = fill_with_linear_regression(real_data)
                    elif answer == "multiple_imputation":
                        print("h")
                        result_data = fill_with_multiple_imputation(real_data)
                    elif answer == "ffill":
                        print("ffff")
                        result_data = fill_with_ffill_bfill(real_data, method="ffill")
                    elif answer == "bfill":
                        print("dasdsadsdash")
                        result_data = fill_with_ffill_bfill(real_data, method="bfill")
                    elif answer == "drop_rows":
                        print("daadsdsadash")
                        result_data = drop_missing(real_data, axis=0)
                    elif answer == "drop_columns":
                        print("hsdadasadssadsaddsadsa")
                        result_data = drop_missing(real_data, axis=1)
                    elif answer == "pchip":
                        print("hsdadasadssadsaddsadsa")
                        result_data= fill_with_pchip(real_data)
                    elif answer == "linear_interpolation":
                        print("hsdadasadssadsaddsadsa")
                        result_data= fill_with_linear(real_data)
                    elif answer == "neighbor_avg":
                        print("hsdadasadssadsaddsadsa")
                        result_data=fill_with_neighbor_average(real_data)
                    elif answer == "mice":
                        print("hsdadasadssadsaddsadsa")
                        result_data=fill_with_mice(real_data) 
                    else:
                        return JsonResponse({"status": "error", "message": f"Unknown answer: {answer}"})
                    # Save the processed data to the output path
                    result_data.to_csv(output_path, index=False)
                    return JsonResponse({"status": "success", "message": f"File saved at {output_path}"})
                except Exception as e:
                    return JsonResponse({"status": "error", "message": str(e)})
            else:
                return JsonResponse({"status": "error", "message": "Invalid input/output paths"})
        except Exception as e:
            return JsonResponse({"status": "error", "message": str(e)})
    return JsonResponse({"status": "error", "message": "Invalid request method"})

# 1. Sabit bir değerle doldurma
def fill_with_constant(data, value=0):
    data.fillna(value, inplace=True)
    return data

# 2. Ortalama/Medyan ile doldurma
def fill_with_mean_or_median(data, method='mean'):
    if method == 'mean':
        data.fillna(data.mean(), inplace=True)
    elif method == 'median':
        data.fillna(data.median(), inplace=True)
    return data

# 3. KNN ile doldurma
def fill_with_knn(data, n_neighbors=5):
    imputer = KNNImputer(n_neighbors=n_neighbors)
    return pd.DataFrame(imputer.fit_transform(data), columns=data.columns)

# 4. Doğrusal Regresyon ile doldurma
# Yöntem 4: Linear Regression ile doldurma
def fill_with_linear_regression(data):
    data = data.copy()  # Orijinal veriyi değiştirmemek için bir kopya oluşturuyoruz
    
    for col in data.columns:
        if data[col].isnull().sum() > 0:  # Eğer sütunda eksik veri varsa
            # Hedef sütun dışındaki tüm sütunlardaki eksik verileri geçici olarak doldur
            temp_data = data.apply(lambda x: x.fillna(x.mean()), axis=0)
            
            X_train = temp_data.loc[data[col].notnull(), data.columns != col]  # Hedef sütunun eksik olmayan değerleri
            y_train = data.loc[data[col].notnull(), col]  # Eksik olmayan değerler
            X_test = temp_data.loc[data[col].isnull(), data.columns != col]  # Eksik olan değerler için özellikler
            
            # Veriyi ölçeklendirme (özellikle farklı boyutlarda veriler için önemli)
            scaler = StandardScaler()
            X_train_scaled = scaler.fit_transform(X_train)
            X_test_scaled = scaler.transform(X_test)
            
            # Linear Regression modelini oluştur ve eğit
            model = LinearRegression()
            model.fit(X_train_scaled, y_train)
            
            # Eksik değerleri modelin tahmin ettiği değerlerle doldur
            data.loc[data[col].isnull(), col] = model.predict(X_test_scaled)
    
    return data


# 5. Çoklu İmputasyon ile doldurma
def fill_with_multiple_imputation(data):
    imputer = IterativeImputer()
    return pd.DataFrame(imputer.fit_transform(data), columns=data.columns)

# 6. İleri veya geri doldurma
def fill_with_ffill_bfill(data, method='ffill'):
    data.fillna(method=method, inplace=True)
    return data

# 7. Eksik veri içeren satırları/sütunları silme
def drop_missing(data, axis=0):
    data.dropna(axis=axis, inplace=True)
    return data


def fill_with_pchip(data):
    data = data.copy()
    for column in data.columns:
        max_value = data[column].max()
        min_value = data[column].min()

        n = data[column].isna().sum()
        if n > 0:
            x = np.arange(len(data))
            y = data[column].values
            mask = ~np.isnan(y)
            
            interpolator = PchipInterpolator(x[mask], y[mask], extrapolate=False)
            y_interp = interpolator(x)
            y_interp = np.clip(y_interp, min_value, max_value)  # Değerleri min/max ile sınırlıyoruz
            
            data[column] = y_interp
    return data    
def fill_with_linear(data):
    data = data.copy()
    for column in data.columns:
        if data[column].isna().sum() > 0:
            data[column] = data[column].interpolate(method='linear', limit_direction='both')
    return data
import pandas as pd

def fill_with_neighbor_average(data):
    data = data.copy()
    for column in data.columns:
        if data[column].isna().sum() > 0:
            for i in range(len(data)):
                # Hücre boşsa doldurma işlemi yap
                if pd.isna(data.loc[i, column]):
                    # Eğer önceki ve sonraki komşular mevcut ve boş değilse ortalamayı al
                    if i > 0 and i < len(data) - 1 and not pd.isna(data.loc[i - 1, column]) and not pd.isna(data.loc[i + 1, column]):
                        interpolated_value = (data.loc[i - 1, column] + data.loc[i + 1, column]) / 2
                        data.loc[i, column] = interpolated_value
                    # Sadece önceki komşu varsa
                    elif i > 0 and not pd.isna(data.loc[i - 1, column]):
                        data.loc[i, column] = data.loc[i - 1, column]
                    # Sadece sonraki komşu varsa
                    elif i < len(data) - 1 and not pd.isna(data.loc[i + 1, column]):
                        data.loc[i, column] = data.loc[i + 1, column]
    return data

def fill_with_mice(data, max_iter=10, random_state=42):
    data = data.copy()
    
    # MICE için iteratif bir imputer oluşturuluyor
    mice_imputer = IterativeImputer(max_iter=max_iter, random_state=random_state)
    
    # Veriye imputation uygulanıyor
    imputed_data = mice_imputer.fit_transform(data)
    
    # İmpute edilen veriyi DataFrame formatına geri çeviriyoruz
    imputed_data_df = pd.DataFrame(imputed_data, columns=data.columns)
    
    return imputed_data_df
