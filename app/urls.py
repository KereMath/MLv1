from django.urls import path
from . import views

urlpatterns = [
    path('fake-deletion/', views.fake_deletion, name='fake_deletion'),
    path('data-imputation/', views.data_imputation, name='data_imputation'),  # Yeni imputation API

]
