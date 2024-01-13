from django.db import models
from ..user.models import User
import uuid


class Profile(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    company_name = models.CharField(max_length=100, null=True)
    logo_image = models.TextField(null=True)
    phone_no = models.IntegerField(null=True, blank=True)
    address = models.CharField(max_length=100, null=True)
    email = models.CharField(max_length=100, null=True)
    city = models.CharField(max_length=100, null=True)
    title = models.CharField(max_length=100, null=True)
    signature_image = models.TextField(null=True)
    manual_signature_image = models.CharField(max_length=20000, null=True)
    user = models.OneToOneField(User, on_delete=models.CASCADE, null=True)
