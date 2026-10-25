class Slot {
  const Slot({
    required this.id,
    required this.provider,
    required this.start,
    required this.end,
    required this.available,
  });

  final String id;
  final String provider;
  final String start;
  final String end;
  final bool available;

  factory Slot.fromJson(Map<String, dynamic> json) => Slot(
    id: json['id'] as String,
    provider: json['provider'] as String,
    start: json['start'] as String,
    end: json['end'] as String,
    available: json['available'] as bool,
  );
}
