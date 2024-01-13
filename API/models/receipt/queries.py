import graphene
from .models import Receipt
from .type import ReceiptType


class ReceiptQuery(graphene.ObjectType):
    all_receipts = graphene.List(ReceiptType, user_id=graphene.UUID(required=True))

    def resolve_all_receipts(self, info, user_id):
        return Receipt.objects.filter(user_id=user_id)
