// lib/pages/login_page.dart
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart'; // token saklamak için

class LoginPage extends StatefulWidget {
  const LoginPage({Key? key}) : super(key: key);
  @override
  _LoginPageState createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _formKey = GlobalKey<FormState>();
  final TextEditingController emailOrUsernameController = TextEditingController();
  final TextEditingController passwordController        = TextEditingController();

  bool isLoading    = false;
  String errorMessage = '';
  bool _obscurePassword = true;

  Future<void> login() async {
  if (!_formKey.currentState!.validate()) return;
  setState(() { isLoading = true; errorMessage = ''; });

  final url = Uri.parse('http://10.0.3.153:5000/login');
  final body = jsonEncode({
    'email'   : emailOrUsernameController.text.trim(),
    'username': emailOrUsernameController.text.trim(),
    'password': passwordController.text.trim(),
  });

  try {
    final response = await http.post(
      url,
      headers: {"Content-Type": "application/json"},
      body: body,
    );
    if (response.statusCode == 200) {
      final rawCookie = response.headers['set-cookie'];
      if (rawCookie == null) {
        setState(() => errorMessage = 'Sunucu cookie döndürmedi.');
        return;
      }
      final token = rawCookie.split(';')[0].split('=')[1];
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('jwt_token', token);

      // Artık sonraki isteklerde:
      // headers: {
      //   "Content-Type":"application/json",
      //   "Cookie":"jwt=$token"
      // }
      Navigator.pushReplacementNamed(context, '/home');
    } else {
      setState(() => errorMessage = 'Giriş başarısız.');
    }
  } catch (err) {
    setState(() => errorMessage = 'Bir hata oluştu: $err');
  } finally {
    setState(() { isLoading = false; });
  }
}


  @override
  void dispose() {
    emailOrUsernameController.dispose();
    passwordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Giriş Yap'), centerTitle: true),
      backgroundColor: Colors.grey[100],
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Card(
            elevation: 4, shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
            child: Padding(
              padding: const EdgeInsets.all(24),
              child: Column(
                children: [
                  const Icon(Icons.login_rounded, size: 80, color: Colors.greenAccent),
                  const SizedBox(height: 10),
                  TextFormField(
                    controller: emailOrUsernameController,
                    decoration: const InputDecoration(labelText: 'E-posta veya Kullanıcı Adı'),
                    validator: (v) => (v==null||v.isEmpty) ? 'Giriş bilgisi giriniz' : null,
                  ),
                  const SizedBox(height: 12),
                  TextFormField(
                    controller: passwordController,
                    obscureText: _obscurePassword,
                    decoration: InputDecoration(
                      labelText: 'Şifre',
                      suffixIcon: IconButton(
                        icon: Icon(_obscurePassword
                            ? Icons.visibility_off
                            : Icons.visibility),
                        onPressed: () => setState(() => _obscurePassword = !_obscurePassword),
                      ),
                    ),
                    validator: (v) => (v==null||v.isEmpty) ? 'Şifre giriniz' : null,
                  ),
                  Align(
                    alignment: Alignment.centerRight,
                    child: TextButton(onPressed: () {}, child: const Text('Şifremi Unuttum?')),
                  ),
                  const SizedBox(height: 8),
                  if (errorMessage.isNotEmpty)
                    Text(errorMessage, style: const TextStyle(color: Colors.red)),
                  const SizedBox(height: 8),
                  isLoading
                      ? const CircularProgressIndicator()
                      : ElevatedButton(
                          onPressed: login,
                          child: const Text('Giriş Yap'),
                          style: ElevatedButton.styleFrom(minimumSize: const Size.fromHeight(50)),
                        ),
                  const SizedBox(height: 16),
                  TextButton(
                    onPressed: () => Navigator.pushReplacementNamed(context, '/signup'),
                    child: const Text('Hesabınız yok mu? Kayıt Olun'),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
