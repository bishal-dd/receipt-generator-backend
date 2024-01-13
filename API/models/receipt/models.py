from django.db import models
from ..user.models import User
import uuid
from django.utils import timezone


class Receipt(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    receipt_name = models.CharField(max_length=100, null=False)
    recipient_name = models.CharField(max_length=100, null=False)
    recipient_phone = models.IntegerField(null=False)
    amount = models.FloatField(null=False)
    journal_no = models.IntegerField(null=True, blank=True)
    date = models.DateField(null=False)
    created_at = models.DateTimeField(default=timezone.now, editable=False)
    updated_at = models.DateTimeField(default=timezone.now, editable=False)
    total_amount = models.FloatField(null=True, default=0)
    user = models.ForeignKey(User, on_delete=models.CASCADE)


