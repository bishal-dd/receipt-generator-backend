import graphene
from .models.user.queries import UserQuery
from .models.receipt.queries import ReceiptQuery
from .models.user.mutations import CreateUserMutation
from .models.receipt.mutations import CreateReceiptMutation


class RootQuery(UserQuery, ReceiptQuery, graphene.ObjectType):
    pass


class RootMutation(graphene.ObjectType):
    create_user = CreateUserMutation.Field()
    create_receipt = CreateReceiptMutation.Field()


schema = graphene.Schema(query=RootQuery, mutation=RootMutation)
