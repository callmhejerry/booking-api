import 'package:flutter/material.dart';

import '../data/booking.dart';

class BookingTile extends StatelessWidget {
  const BookingTile({
    super.key,
    required this.booking,
    required this.cancelling,
    required this.onCancel,
  });

  final Booking booking;
  final bool cancelling;
  final VoidCallback? onCancel;

  @override
  Widget build(BuildContext context) {
    final statusColor = booking.status.toLowerCase() == 'cancelled'
        ? const Color(0xFF79847F)
        : const Color(0xFF176B55);

    return Container(
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: const Color(0xFFFAFBF8),
        borderRadius: BorderRadius.circular(14),
        border: Border.all(color: const Color(0xFFE7ECE6)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(Icons.event_note, color: Color(0xFF167D73)),
              const SizedBox(width: 10),
              Expanded(
                child: Text(
                  booking.name,
                  style: const TextStyle(fontWeight: FontWeight.w700),
                ),
              ),
              Text(
                booking.status,
                style: TextStyle(
                  color: statusColor,
                  fontWeight: FontWeight.w700,
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          SelectableText('Booking ID: ${booking.id}'),
          Text('Slot: ${booking.slotId}'),
          Text(booking.email),
          Text('Created: ${booking.createdAt}'),
          if (onCancel != null) ...[
            const SizedBox(height: 8),
            Align(
              alignment: Alignment.centerRight,
              child: TextButton.icon(
                onPressed: cancelling ? null : onCancel,
                icon: const Icon(Icons.event_busy),
                label: Text(cancelling ? 'Cancelling…' : 'Cancel booking'),
                style: TextButton.styleFrom(
                  foregroundColor: const Color(0xFF9B2929),
                ),
              ),
            ),
          ],
        ],
      ),
    );
  }
}
