import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:jwt_decoder/jwt_decoder.dart';
import 'package:qrfatura/screen/models.dart';
import 'package:qrfatura/screen/scan.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:intl/intl.dart';

class MainPage extends StatefulWidget {
  const MainPage({Key? key}) : super(key: key);

  @override
  State<MainPage> createState() => _MainPageState();
}

class _MainPageState extends State<MainPage> {
  List<Invoice> _invoices = [];
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _fetchInvoices();
  }

  Future<void> _fetchInvoices() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final token = prefs.getString('jwt_token');
      if (token == null) throw 'Token bulunamadı';

      final claims = JwtDecoder.decode(token);
      final userId = claims['sub'].toString();

      final url = Uri.parse('http://10.0.3.153:5000/invoices/user/$userId');
      final res = await http.get(
        url,
        headers: {
          'Content-Type': 'application/json',
          'Cookie': 'Auth=$token',
        },
      );

      if (res.statusCode == 200) {
        final body = jsonDecode(res.body);
        final List<dynamic> rawList = body is List ? body : [body];
        setState(() {
          _invoices = rawList
              .map((e) => Invoice.fromJson(e as Map<String, dynamic>))
              .toList();
          _loading = false;
        });
      } else {
        throw 'Sunucu hatası: ${res.statusCode}';
      }
    } catch (e) {
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  } // <-- Burada _fetchInvoices kapandı

  Future<void> _assignInvoice(String invoiceId) async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('jwt_token')!;
    final userId = JwtDecoder.decode(token)['sub'].toString();

    final url = Uri.parse('http://10.0.3.50:5000/invoices/$invoiceId/assign');
    final res = await http.put(
      url,
      headers: {
        'Content-Type': 'application/json',
        'Cookie': 'Auth=$token',
      },
      body: jsonEncode({'user_id': int.parse(userId)}),
    );

    if (res.statusCode == 200) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Fatura kaydedildi!')),
      );
      await _fetchInvoices();
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Hata: ${res.statusCode}')),
      );
    }
  }

  String _formatDate(DateTime date) =>
      DateFormat('dd.MM.yyyy').format(date);

  String _formatCurrency(double amount) =>
      '₺${amount.toStringAsFixed(2)}';

  @override
  Widget build(BuildContext context) {
    if (_loading) {
      return const Scaffold(
        body: Center(child: CircularProgressIndicator()),
      );
    }
    if (_error != null) {
      return Scaffold(
        body: Center(child: Text('Hata: $_error')),
      );
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text('Faturalarım'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () {
              setState(() => _loading = true);
              _fetchInvoices();
            },
          ),
        ],
      ),
      drawer: Drawer(
        child: ListView(
          padding: EdgeInsets.zero,
          children: [
            const DrawerHeader(
              decoration: BoxDecoration(color: Colors.blue),
              child: Text(
                'Menü',
                style: TextStyle(color: Colors.white, fontSize: 24),
              ),
            ),
            ListTile(
              leading: const Icon(Icons.qr_code_scanner),
              title: const Text('QR Kodu Okut'),
              onTap: () async {
                Navigator.of(context).pop();
                final scannedUrl = await Navigator.of(context).push<String>(
                  MaterialPageRoute(builder: (_) => const ScannerPage()),
                );
                if (scannedUrl != null) {
                  try {
                    final uri = Uri.parse(scannedUrl);
                    final segments = uri.pathSegments;
                    if (segments.length >= 2 &&
                        segments[segments.length - 2] == 'invoices') {
                      final invId = segments.last;
                      await _assignInvoice(invId);
                    } else {
                      throw 'Geçersiz URL';
                    }
                  } catch (e) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text('Okuma hatası: $e')),
                    );
                  }
                }
              },
            ),
          ],
        ),
      ),
      body: _invoices.isEmpty
          ? const Center(child: Text('Henüz faturanız bulunmamaktadır.'))
          : ListView.builder(
              itemCount: _invoices.length,
              itemBuilder: (ctx, i) {
                final inv = _invoices[i];
                return Card(
                  margin:
                      const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                  elevation: 3,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: ExpansionTile(
                    tilePadding: const EdgeInsets.symmetric(
                        horizontal: 16, vertical: 8),
                    childrenPadding:
                        const EdgeInsets.fromLTRB(16, 0, 16, 16),
                    expandedCrossAxisAlignment: CrossAxisAlignment.start,
                    title: Text(
                      '${inv.place} mağazasında ${_formatDate(inv.createdAt)} tarihli alışveriş',
                      style: const TextStyle(fontWeight: FontWeight.bold),
                    ),
                    subtitle: Text('Toplam: ${_formatCurrency(inv.total)}'),
                    children: [
                      const Divider(),
                      const SizedBox(height: 8),
                      Text(
                        'Alışveriş Yeri: ${inv.place}',
                        style: const TextStyle(fontSize: 16),
                      ),
                      const SizedBox(height: 12),
                      const Text(
                        'Ürünler',
                        style: TextStyle(
                            fontWeight: FontWeight.bold, fontSize: 16),
                      ),
                      const SizedBox(height: 8),
                      Table(
                        columnWidths: const {
                          0: FlexColumnWidth(3),
                          1: FlexColumnWidth(1),
                          2: FlexColumnWidth(2),
                          3: FlexColumnWidth(2),
                        },
                        border: TableBorder.all(
                          color: Colors.grey.shade300,
                          width: 1,
                        ),
                        children: [
                          TableRow(
                            decoration: BoxDecoration(
                              color: Colors.grey.shade200,
                            ),
                            children: const [
                              Padding(
                                padding: EdgeInsets.all(8.0),
                                child: Text('Ürün Adı',
                                    style:
                                        TextStyle(fontWeight: FontWeight.bold)),
                              ),
                              Padding(
                                padding: EdgeInsets.all(8.0),
                                child: Text('Adet',
                                    style:
                                        TextStyle(fontWeight: FontWeight.bold)),
                              ),
                              Padding(
                                padding: EdgeInsets.all(8.0),
                                child: Text('Birim Fiyat',
                                    style:
                                        TextStyle(fontWeight: FontWeight.bold)),
                              ),
                              Padding(
                                padding: EdgeInsets.all(8.0),
                                child: Text('Ara Toplam',
                                    style:
                                        TextStyle(fontWeight: FontWeight.bold)),
                              ),
                            ],
                          ),
                          ...inv.items.map(
                            (item) => TableRow(
                              children: [
                                Padding(
                                  padding: const EdgeInsets.all(8.0),
                                  child: Text(item.productName),
                                ),
                                Padding(
                                  padding: const EdgeInsets.all(8.0),
                                  child: Text('${item.quantity}'),
                                ),
                                Padding(
                                  padding: const EdgeInsets.all(8.0),
                                  child: Text('₺${item.unitPrice}'),
                                ),
                                Padding(
                                  padding: const EdgeInsets.all(8.0),
                                  child:
                                      Text(_formatCurrency(item.subtotal)),
                                ),
                              ],
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 16),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const Text('KDV Dahil Toplam:'),
                          Text(
                            _formatCurrency(inv.total),
                            style: const TextStyle(
                                fontWeight: FontWeight.bold, fontSize: 16),
                          ),
                        ],
                      ),
                      const SizedBox(height: 16),
                      Text('Kullanıcı: ${inv.username}'),
                      if (inv.email.isNotEmpty)
                        Text('E-posta: ${inv.email}'),
                    ],
                  ),
                );
              },
            ),
    );
  }
}
