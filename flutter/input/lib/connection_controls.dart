import 'package:flutter/material.dart';
import 'package:input/app_state.dart';
import 'package:input/bloc.dart';

class ConnectionControls extends StatefulWidget {
  createState() => _ConnectionControlsState();
}

class _ConnectionControlsState extends State<ConnectionControls> {
  AppState state;
  Future futurePassword;
  TextEditingController urlController;
  TextEditingController userController;
  TextEditingController passwordController;
  TextEditingController databaseController;

  initState(){
    state = BlocProvider.of<AppState>(context);
    futurePassword = getPassword();
    urlController = TextEditingController(text: state.connStr);
    userController = TextEditingController(text: state.username);
    databaseController = TextEditingController(text: state.database);
    super.initState();
  }

  Future getPassword() async {
    var password = await state.getPassword();
    passwordController = TextEditingController(text: password);
  }

  void urlChanged(value) {
    state.connStr = value;
  }

  void usernameChanged(value) {
    state.username = value;
  }

  void passwordChanged(value) {
    state.password = value;
  }

  void databaseChanged(value) {
    state.database = value;
  }

  void toggleUrlCollapse() {
    setState(() {
      state.urlCollapsed = !state.urlCollapsed;
    });
  }

  Widget build(BuildContext context) {
    return Card(
      elevation: 12,
      child: FutureBuilder(
        future: futurePassword,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return Center(child: CircularProgressIndicator());
          }
          return Padding(
            padding: EdgeInsets.fromLTRB(8, 0, 8, 8),
            child: state.urlCollapsed ? Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                SizedBox(height: 16, child: TextButton(onPressed: toggleUrlCollapse, child: Icon(Icons.keyboard_arrow_down))),
                Text(this.urlController.text, overflow: TextOverflow.ellipsis),
              ],
            ) : Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                SizedBox(height: 16, child: TextButton(onPressed: toggleUrlCollapse, child: Icon(Icons.keyboard_arrow_up))),
                TextField(controller: this.urlController, decoration: InputDecoration(labelText: "connection (bolt or neo4j address)"), onChanged: urlChanged, autocorrect: false),
                TextField(controller: this.userController, decoration: InputDecoration(labelText: "username"), onChanged: usernameChanged, autocorrect: false),
                TextField(controller: this.passwordController, decoration: InputDecoration(labelText: "password"), autocorrect: false, obscureText: true, onChanged: passwordChanged),
                TextField(controller: this.databaseController, decoration: InputDecoration(labelText: "database"), onChanged: databaseChanged, autocorrect: false),
              ],
            ),
          );
        },
      ),
    );
  }
}
