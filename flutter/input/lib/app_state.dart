
import 'package:input/bloc.dart';
import 'package:rxdart/rxdart.dart';

class AppState extends BlocState {
  var _isLoaded = BehaviorSubject<bool>.seeded(false);
  Stream<bool> get isLoaded => _isLoaded.stream;

  void dispose() {
    _isLoaded.close();
  }

  Future initialize() {
  }
}