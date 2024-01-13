import graphene
from .models import User
from .type import UserType


class UserQuery(graphene.ObjectType):
    all_users = graphene.List(UserType)

    def resolve_all_users(root, info):
        return User.objects.all()
