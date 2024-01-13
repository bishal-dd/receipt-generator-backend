import graphene
from graphene_django import DjangoObjectType

from .models import Profile
from ..user.type import UserType


class ProfileType(DjangoObjectType):
    user = graphene.Field(UserType)

    class Meta:
        model = Profile
        fields = ("id", "company_name", "logo_image", "phone_no",
                  "address", "email", "city", "title", "signature_image",
                  "manual_signature_image", "user")

    def resolve_user(self, info):
        return self.user
