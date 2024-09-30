from django.urls import path
from . import views

urlpatterns = [
    path('fake-deletion/', views.fake_deletion, name='fake_deletion'),
]
