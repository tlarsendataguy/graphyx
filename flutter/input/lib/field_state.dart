import 'package:input/bloc.dart';
import 'package:rxdart/rxdart.dart' as rx;
import 'package:input/configuration.dart' as c;

class FieldState extends BlocState {
  FieldState(this.field);
  c.Field field;

  var _pathChanged = rx.PublishSubject();
  Stream get pathChanged => _pathChanged.stream;

  void addElementToPath(c.Element element){
    field.DataType = element.DataType;
    field.Path.add(element);
    _pathChanged.add(null);
  }

  void dispose() {
    _pathChanged.close();
  }

  Future initialize() {
  }

}