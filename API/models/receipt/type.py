import graphene
from graphene_django import DjangoObjectType

from .models import Receipt
from ..user.type import UserType


class ReceiptType(DjangoObjectType):
    user = graphene.Field(UserType)

    class Meta:
        model = Receipt
        fields = ("id", "receipt_name", "recipient_name",
                  "recipient_phone", "amount", "journal_no",
                  "date", "created_at", "updated_at",
                  "total_amount", "user")

    def resolve_user(self, info):
        return self.user


class ReceiptInputType(graphene.InputObjectType):
    user_id = graphene.UUID(required=True)
    receipt_name = graphene.String(required=True)
    recipient_name = graphene.String(required=True)
    recipient_phone = graphene.String(required=True)
    amount = graphene.Int(required=True)
    journal_no = graphene.Int()
    total_amount = graphene.Int()
    date = graphene.Date()
