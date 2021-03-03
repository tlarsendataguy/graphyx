import 'package:input/bloc.dart';
import 'package:rxdart/rxdart.dart' as rx;
import 'package:input/configuration.dart' as c;

class FieldState extends BlocState {
  FieldState(this.field);
  c.FieldData field;

  var _pathChanged = rx.PublishSubject();
  Stream get pathChanged => _pathChanged.stream;

  void addElementToPath(c.ElementData element){
    field.DataType = element.DataType;
    field.Path.add(c.ElementContainer(Element: element));
    _pathChanged.add(null);
  }

  void truncatePathAtElement(int index) {
    field.Path = field.Path.sublist(0, index);
    _pathChanged.add(null);
  }

  void dispose() {
    _pathChanged.close();
  }

  Future initialize() {
  }

}