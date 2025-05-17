import 'package:flutter/material.dart';
import 'package:qrfatura/screen/homescreen.dart';
import 'package:qrfatura/screen/login.dart';
import 'package:qrfatura/screen/signup.dart';

void main() {
  runApp(Routes());
}

class Routes extends StatelessWidget {
  @override Widget build(BuildContext context) {
    return MaterialApp(
      initialRoute: '/login',
      routes: {
        '/login': (context) => const LoginPage(),
        '/signup': (context) => const SignupPage(),
        '/home': (context) => const MainPage()
      },
    );
  }
}
