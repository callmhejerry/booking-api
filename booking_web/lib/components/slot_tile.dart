import 'package:booking_web/data/slot.dart';
import 'package:flutter/material.dart';

class SlotTile extends StatelessWidget {
  const SlotTile({
    super.key,
    required this.slot,
    required this.selected,
    required this.onTap,
  });

  final Slot slot;
  final bool selected;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final start = DateTime.tryParse(slot.start)?.toLocal();
    final end = DateTime.tryParse(slot.end)?.toLocal();
    final period = start == null
        ? slot.start
        : '${_clock(start)}${end == null ? '' : ' – ${_clock(end)}'}';

    return Material(
      color: selected ? const Color(0xFFE7F4EF) : const Color(0xFFFAFBF8),
      borderRadius: BorderRadius.circular(14),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(14),
        child: Container(
          padding: const EdgeInsets.all(14),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(14),
            border: Border.all(
              color: selected
                  ? const Color(0xFF167D73)
                  : const Color(0xFFE7ECE6),
            ),
          ),
          child: Row(
            children: [
              Icon(
                selected ? Icons.radio_button_checked : Icons.schedule,
                color: selected
                    ? const Color(0xFF167D73)
                    : const Color(0xFF60716A),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      period,
                      style: const TextStyle(fontWeight: FontWeight.w700),
                    ),
                    Text(
                      '${slot.provider} · ${slot.id}',
                      style: const TextStyle(color: Color(0xFF64736D)),
                    ),
                  ],
                ),
              ),
              Text(
                slot.available ? 'Open' : 'Booked',
                style: TextStyle(
                  color: slot.available
                      ? const Color(0xFF176B55)
                      : const Color(0xFF79847F),
                  fontWeight: FontWeight.w600,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  static String _clock(DateTime value) {
    final hour = value.hour % 12 == 0 ? 12 : value.hour % 12;
    final minute = value.minute.toString().padLeft(2, '0');
    return '$hour:$minute ${value.hour < 12 ? 'AM' : 'PM'}';
  }
}
