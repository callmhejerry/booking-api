class Booking {
  const Booking({
    required this.id,
    required this.slotId,
    required this.name,
    required this.email,
    required this.status,
    required this.createdAt,
  });

  final String id;
  final String slotId;
  final String name;
  final String email;
  final String status;
  final String createdAt;

  factory Booking.fromJson(Map<String, dynamic> json) => Booking(
    id: json['id'] as String,
    slotId: json['slot_id'] as String,
    name: json['name'] as String,
    email: json['email'] as String,
    status: json['status'] as String,
    createdAt: json['created_at'] as String,
  );
}
