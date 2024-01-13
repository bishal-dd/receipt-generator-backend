import graphene
from .type import ProfileType
from .models import Profile


class ProfileQuery(graphene.ObjectType):
    profile = graphene.List(ProfileType, user_id=graphene.UUID(required=True))

    def resolve_profile(self, info, user_id):
        return Profile.objects.filter(user_id=user_id)