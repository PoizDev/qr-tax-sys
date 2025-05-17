import 'package:flutter/material.dart';
import 'package:mobile_scanner/mobile_scanner.dart';
import 'package:permission_handler/permission_handler.dart';

class ScannerPage extends StatefulWidget {
  const ScannerPage({super.key});

  @override
  State<ScannerPage> createState() => _ScannerPageState();
}

class _ScannerPageState extends State<ScannerPage> {
  bool _scanned = false;
  bool _hasPermission = false;
  final MobileScannerController _controller = MobileScannerController();

  @override
  void initState() {
    super.initState();
    _requestCameraPermission();
  }

  Future<void> _requestCameraPermission() async {
    final status = await Permission.camera.request();
    setState(() {
      _hasPermission = status.isGranted;
    });

    if (status.isPermanentlyDenied) {
      // Kullanıcı izni kalıcı olarak reddetti, ayarlardan açmaları gerekiyor
      _showPermissionDialog();
    }
  }

  void _showPermissionDialog() {
    showDialog(
      context: context,
      builder: (BuildContext context) => AlertDialog(
        title: const Text('Kamera İzni Gerekli'),
        content: const Text(
            'QR kod taraması için kamera izni gereklidir. Lütfen ayarlardan kamera iznini etkinleştirin.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('İptal'),
          ),
          TextButton(
            onPressed: () {
              openAppSettings();
              Navigator.of(context).pop();
            },
            child: const Text('Ayarları Aç'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('QR Kodu Okut'),
        actions: [
          IconButton(
            icon: const Icon(Icons.flash_on),
            onPressed: _hasPermission ? () => _controller.toggleTorch() : null,
          ),
          IconButton(
            icon: const Icon(Icons.flip_camera_ios),
            onPressed: _hasPermission ? () => _controller.switchCamera() : null,
          ),
        ],
      ),
      body: !_hasPermission
          ? Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Text(
                    'Kamera izni verilmedi',
                    style: TextStyle(fontSize: 18),
                  ),
                  const SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: _requestCameraPermission,
                    child: const Text('İzin İste'),
                  ),
                ],
              ),
            )
          : MobileScanner(
              controller: _controller,
              onDetect: (capture) {
                final List<Barcode> barcodes = capture.barcodes;
                // Check if we have any valid barcodes
                if (barcodes.isNotEmpty && !_scanned) {
                  for (final barcode in barcodes) {
                    // Make sure we have a value before proceeding
                    if (barcode.rawValue != null) {
                      setState(() {
                        _scanned = true;
                      });
                      // Return to previous screen with the scanned value
                      Navigator.of(context).pop(barcode.rawValue);
                      break; // We only need the first valid code
                    }
                  }
                }
              },
            ),
    );
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }
}