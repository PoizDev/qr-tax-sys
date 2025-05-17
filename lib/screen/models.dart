class InvoiceItem {
  final int id;
  final int productId;
  final String productName;
  final String unitPrice;
  final double taxRate;
  final int quantity;
  final double subtotal;

  InvoiceItem({
    required this.id,
    required this.productId,
    required this.productName,
    required this.unitPrice,
    required this.taxRate,
    required this.quantity,
    required this.subtotal,
  });

  factory InvoiceItem.fromJson(Map<String, dynamic> json) {
    // Parse product information from nested product object
    final product = json['product'] as Map<String, dynamic>?;
    
    return InvoiceItem(
      id: json['id'],
      productId: json['product_id'],
      productName: product != null ? product['product_name'] : 'Ürün Bilgisi Yok',
      unitPrice: json['unit_price'],
      taxRate: (json['tax_rate'] as num).toDouble(),
      quantity: json['quantity'],
      subtotal: _calculateSubtotal(
        json['quantity'], 
        double.tryParse(json['unit_price']) ?? 0.0, 
        (json['tax_rate'] as num).toDouble()
      ),
    );
  }

  static double _calculateSubtotal(int quantity, double unitPrice, double taxRate) {
    return quantity * unitPrice * (1 + taxRate);
  }
}

class Invoice {
  final int id;
  final List<InvoiceItem> items;
  final double total;
  final String place;
  final DateTime createdAt;
  final String username;
  final String email;

  Invoice({
    required this.id,
    required this.items,
    required this.total,
    required this.place,
    required this.createdAt,
    required this.username,
    required this.email,
  });

  factory Invoice.fromJson(Map<String, dynamic> json) {
    var itemsList = (json['items'] as List)
        .map((i) => InvoiceItem.fromJson(i))
        .toList();
    
    final user = json['user'] as Map<String, dynamic>?;
    final username = user != null ? user['username'] : 'Bilinmeyen Kullanıcı';
    final email = user != null ? user['email'] : '';
    
    return Invoice(
      id: json['id'],
      items: itemsList,
      total: (json['total'] as num).toDouble(),
      place: json['place'],
      createdAt: DateTime.parse(json['created_at']),
      username: username,
      email: email,
    );
  }
}