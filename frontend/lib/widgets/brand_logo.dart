import 'package:flutter/material.dart';

/// Renders the official ToolGridLab vector brand mark (Connected 'T' Modular Grid + Active Lab Node).
class BrandLogo extends StatelessWidget {
  final double size;
  final bool showGlow;

  const BrandLogo({
    super.key,
    this.size = 38,
    this.showGlow = true,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      width: size,
      height: size,
      decoration: showGlow
          ? BoxDecoration(
              borderRadius: BorderRadius.circular(size * 0.22),
              boxShadow: [
                BoxShadow(
                  color: const Color(0xFF3B82F6).withValues(alpha: 0.28),
                  blurRadius: size * 0.32,
                  offset: Offset(0, size * 0.1),
                ),
              ],
            )
          : null,
      child: CustomPaint(
        size: Size.square(size),
        painter: _BrandLogoPainter(),
      ),
    );
  }
}

class _BrandLogoPainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final s = size.width / 512.0;

    RRect rrect(double x, double y, double w, double h, double r) {
      return RRect.fromRectAndRadius(
        Rect.fromLTWH(x * s, y * s, w * s, h * s),
        Radius.circular(r * s),
      );
    }

    // 1. Outer Squircle Base
    final outerRect = Rect.fromLTWH(16 * s, 16 * s, 480 * s, 480 * s);
    final outerRRect = RRect.fromRectAndRadius(outerRect, Radius.circular(112 * s));

    final bgPaint = Paint()
      ..shader = const LinearGradient(
        begin: Alignment.topLeft,
        end: Alignment.bottomRight,
        colors: [
          Color(0xFF0B1120),
          Color(0xFF0F172A),
          Color(0xFF1E1B4B),
        ],
        stops: [0.0, 0.55, 1.0],
      ).createShader(outerRect);
    canvas.drawRRect(outerRRect, bgPaint);

    // Outer Border Glow
    final borderPaint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 12 * s
      ..shader = const LinearGradient(
        begin: Alignment.topLeft,
        end: Alignment.bottomRight,
        colors: [
          Color(0xD93B82F6),
          Color(0x806366F1),
          Color(0xBF10B981),
        ],
      ).createShader(outerRect);
    canvas.drawRRect(outerRRect, borderPaint);

    // 2. Inactive Modular Grid Slots
    final slotFillPaint = Paint()..color = const Color(0xFF1E293B);
    final slotStrokePaint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 6 * s
      ..color = const Color(0xFF334155);

    const inactiveSlots = [
      [96.0, 212.0],
      [328.0, 212.0],
      [96.0, 328.0],
    ];
    for (final pos in inactiveSlots) {
      final rr = rrect(pos[0], pos[1], 88, 88, 22);
      canvas.drawRRect(rr, slotFillPaint);
      canvas.drawRRect(rr, slotStrokePaint);
    }

    // 3. Primary "T" Core Architecture (Bridges + 5 Tiles)
    final tBounds = Rect.fromLTWH(96 * s, 96 * s, 320 * s, 320 * s);
    final tPaint = Paint()
      ..shader = const LinearGradient(
        begin: Alignment.topLeft,
        end: Alignment.bottomRight,
        colors: [
          Color(0xFF60A5FA),
          Color(0xFF3B82F6),
          Color(0xFF8B5CF6),
        ],
        stops: [0.0, 0.45, 1.0],
      ).createShader(tBounds);

    const bridges = [
      [176.0, 124.0, 44.0, 32.0],
      [292.0, 124.0, 44.0, 32.0],
      [240.0, 176.0, 32.0, 44.0],
      [240.0, 292.0, 32.0, 44.0],
    ];
    for (final b in bridges) {
      canvas.drawRRect(rrect(b[0], b[1], b[2], b[3], 8), tPaint);
    }

    const tTiles = [
      [96.0, 96.0],
      [212.0, 96.0],
      [328.0, 96.0],
      [212.0, 212.0],
      [212.0, 328.0],
    ];
    for (final pos in tTiles) {
      canvas.drawRRect(rrect(pos[0], pos[1], 88, 88, 22), tPaint);
    }

    // 4. Active Lab Emerald Micro-Tool Node (Bottom-Right)
    final emeraldRect = Rect.fromLTWH(328 * s, 328 * s, 88 * s, 88 * s);
    final emeraldPaint = Paint()
      ..shader = const LinearGradient(
        begin: Alignment.topLeft,
        end: Alignment.bottomRight,
        colors: [
          Color(0xFF34D399),
          Color(0xFF059669),
        ],
      ).createShader(emeraldRect);
    canvas.drawRRect(rrect(328, 328, 88, 88, 22), emeraldPaint);

    // Core dot inside active node
    final dotPaint = Paint()..color = const Color(0xFFECFDF5);
    canvas.drawCircle(Offset(372 * s, 372 * s), 16 * s, dotPaint);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}

