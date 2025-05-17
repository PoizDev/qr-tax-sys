import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

class SignupPage extends StatefulWidget {
  const SignupPage({Key? key}) : super(key: key);
  @override
  _SignupPageState createState() => _SignupPageState();
}

class _SignupPageState extends State<SignupPage> {
  final _formKey = GlobalKey<FormState>();
  final TextEditingController emailController    = TextEditingController();
  final TextEditingController usernameController = TextEditingController();
  final TextEditingController passwordController = TextEditingController();

  bool isLoading    = false;
  String errorMessage = '';

  Future<void> signup() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() { isLoading = true; errorMessage = ''; });

    final url  = Uri.parse('http://10.0.3.153:5000/signup');
    final body = jsonEncode({
      'email'   : emailController.text.trim(),
      'username': usernameController.text.trim(),
      'password': passwordController.text.trim(),
    });

    try {
      final response = await http
          .post(url, headers: {"Content-Type": "application/json"}, body: body)
          .timeout(const Duration(seconds: 10));

      if (response.statusCode == 200) {
        Navigator.pushReplacementNamed(context, '/login');
      } else {
        setState(() {
          errorMessage = 'Kayıt başarısız: ${response.statusCode}';
        });
      }
    } catch (err) {
      setState(() {
        errorMessage = 'Bir hata oluştu: $err';
      });
    } finally {
      setState(() { isLoading = false; });
    }
  }

  @override
  void dispose() {
    emailController.dispose();
    usernameController.dispose();
    passwordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Kayıt Ol'), centerTitle: true),
      backgroundColor: Colors.grey[100],
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Card(
            elevation: 4, shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                children: [
                  const Icon(Icons.person_add_rounded, size: 80, color: Colors.redAccent),
                  const SizedBox(height: 10),
                  TextFormField(
                    controller: emailController,
                    decoration: const InputDecoration(labelText: 'E-posta'),
                    validator: (v) => (v==null||v.isEmpty) ? 'E-posta giriniz' : null,
                  ),
                  const SizedBox(height: 12),
                  TextFormField(
                    controller: usernameController,
                    decoration: const InputDecoration(labelText: 'Kullanıcı Adı'),
                    validator: (v) => (v==null||v.isEmpty) ? 'Kullanıcı adı giriniz' : null,
                  ),
                  const SizedBox(height: 12),
                  TextFormField(
                    controller: passwordController,
                    decoration: const InputDecoration(labelText: 'Şifre'),
                    obscureText: true,
                    validator: (v) => (v==null||v.isEmpty) ? 'Şifre giriniz' : null,
                  ),
                  const SizedBox(height: 20),
                  if (errorMessage.isNotEmpty)
                    Text(errorMessage, style: const TextStyle(color: Colors.red)),
                  const SizedBox(height: 8),
                  isLoading
                      ? const CircularProgressIndicator()
                      : ElevatedButton(
                          onPressed: signup,
                          child: const Text('Kayıt Ol'),
                          style: ElevatedButton.styleFrom(minimumSize: const Size.fromHeight(50)),
                        ),
                  const SizedBox(height: 16),
                  TextButton(
                    onPressed: () => Navigator.pushReplacementNamed(context, '/login'),
                    child: const Text('Zaten hesabınız var mı? Giriş Yapın'),
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
