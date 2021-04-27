import 'package:flutter/widgets.dart';

abstract class BlocState{
  void dispose();
  Future initialize();
}

class BlocProvider<T extends BlocState> extends StatefulWidget{
  BlocProvider({Key key, this.child, this.bloc}) : super(key: key);

  final Widget child;
  final T bloc;

  _BlocProviderState<T> createState() => _BlocProviderState<T>();

  static T of<T extends BlocState>(BuildContext context) {
    BlocProvider<T> provider =
    context.findAncestorWidgetOfExactType<BlocProvider<T>>();
    return provider.bloc;
  }
}

class _BlocProviderState<T extends BlocState> extends State<BlocProvider> {
  void dispose() {
    widget.bloc.dispose();
    super.dispose();
  }

  Widget build(BuildContext context) {
    return widget.child;
  }
}