from .type import UserInputType, UserType
from .models import User
import graphene


class CreateUserMutation(graphene.Mutation):
    user = graphene.Field(UserType)

    class Arguments:
        input_data = UserInputType(required=True)

    def mutate(self, info, input_data):
        user_instance = User(**input_data)
        user_instance.save()

        return CreateUserMutation(user=user_instance)
