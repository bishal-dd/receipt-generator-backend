from graphene_django import DjangoObjectType
import graphene

from .models import User


class UserType(DjangoObjectType):
    class Meta:
        model = User
        fields = ("id", "password", "last_login", "is_superuser", "username",
                  "first_name", "last_name", "email", "is_staff", "is_active", "date_joined")


class UserInputType(graphene.InputObjectType):
    username = graphene.String(required=True)
    email = graphene.String(required=True)
