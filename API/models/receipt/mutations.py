import graphene
from .type import ReceiptType, ReceiptInputType
from .models import Receipt


class CreateReceiptMutation(graphene.Mutation):
    receipt = graphene.Field(ReceiptType)

    class Arguments:
        input_data = ReceiptInputType(required=True)

    def mutate(self, info, input_data):
        receipt_instance = Receipt(**input_data)
        receipt_instance.save()

        return CreateReceiptMutation(receipt=receipt_instance)