
import 'package:input/bloc.dart';
import 'package:rxdart/rxdart.dart' as rx;

enum updated {
  ReturnValues,
  Fields,
}

class AppState extends BlocState {
  var _returnValues = rx.PublishSubject();
  Stream get returnValues => _returnValues.stream;
  var _fields = rx.PublishSubject();
  Stream get fields => _fields.stream;

  void notifyUpdated(updated item) {
    switch (item) {
      case updated.ReturnValues:
        _returnValues.add(null);
        return;
      case updated.Fields:
        _fields.add(null);
        return;
    }
  }

  void dispose() {
    _returnValues.close();
    _fields.close();
  }

  Future initialize() {}
}