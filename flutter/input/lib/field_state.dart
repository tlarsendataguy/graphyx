import 'package:input/app_state.dart';
import 'package:input/bloc.dart';
import 'package:rxdart/rxdart.dart' as rx;

class FieldState extends BlocState {
  FieldState(this.field) {
    _pathChanged = rx.BehaviorSubject<List<PathElement>>.seeded(field.path);
  }
  Field field;

  rx.BehaviorSubject<List<PathElement>> _pathChanged;
  Stream get pathChanged => _pathChanged.stream;

  void addElementToPath(PathElement element){
    field.dataType = element.dataType;
    field.path.add(element);
    _pathChanged.add(field.path);
  }

  void truncatePathAtElement(int index) {
    field.path = field.path.sublist(0, index);
    _pathChanged.add(field.path);
  }

  void dispose() {
    _pathChanged.close();
  }

  Future initialize() {
  }

}