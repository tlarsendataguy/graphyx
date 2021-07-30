import 'package:input/validated_response.dart';

abstract class Validator {
  Future<ValidatedResponse> validate();
}
