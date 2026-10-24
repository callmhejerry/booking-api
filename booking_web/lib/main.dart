import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

void main() {
  runApp(const BookingWebApp());
}

class BookingWebApp extends StatelessWidget {
  const BookingWebApp({super.key});

  @override
  Widget build(BuildContext context) {
    const ink = Color(0xFF182B2A);
    const teal = Color(0xFF167D73);

    return MaterialApp(
      title: 'Book a session',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        useMaterial3: true,
        colorScheme: ColorScheme.fromSeed(
          seedColor: teal,
          surface: const Color(0xFFF7F8F4),
        ),
        scaffoldBackgroundColor: const Color(0xFFF7F8F4),
        appBarTheme: const AppBarTheme(
          backgroundColor: Color(0xFFF7F8F4),
          foregroundColor: ink,
          centerTitle: false,
        ),
        inputDecorationTheme: InputDecorationTheme(
          filled: true,
          fillColor: Colors.white,
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(14),
            borderSide: const BorderSide(color: Color(0xFFDCE4DE)),
          ),
          enabledBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(14),
            borderSide: const BorderSide(color: Color(0xFFDCE4DE)),
          ),
        ),
      ),
      home: const BookingPage(),
    );
  }
}

class BookingPage extends StatefulWidget {
  const BookingPage({super.key});

  @override
  State<BookingPage> createState() => _BookingPageState();
}

class _BookingPageState extends State<BookingPage> {
  final _apiUrlController = TextEditingController();
  final _apiKeyController = TextEditingController();
  final _providerController = TextEditingController();
  final _nameController = TextEditingController();
  final _emailController = TextEditingController();
  final _bookingEmailController = TextEditingController();
  final _bookingIdController = TextEditingController();
  final _resetSlotsController = TextEditingController(
    text: '''
{
  "slots": [
    {
      "id": "slot-1",
      "provider": "Dr Example",
      "start": "2026-12-01T09:00:00Z",
      "end": "2026-12-01T09:30:00Z"
    }
  ]
}''',
  );
  final _client = http.Client();

  List<Slot> _slots = [];
  List<Booking> _bookings = [];
  Slot? _selectedSlot;
  String? _nextBookingCursor;
  bool _loadingSlots = false;
  bool _loadingBookings = false;
  bool _loadingBookingById = false;
  bool _cancellingBooking = false;
  bool _submitting = false;
  bool _resettingSlots = false;
  String? _error;
  String? _message;

  @override
  void dispose() {
    _apiUrlController.dispose();
    _apiKeyController.dispose();
    _providerController.dispose();
    _nameController.dispose();
    _emailController.dispose();
    _bookingEmailController.dispose();
    _bookingIdController.dispose();
    _resetSlotsController.dispose();
    _client.close();
    super.dispose();
  }

  Future<void> _loadSlots() async {
    final baseUrl = _apiUrlController.text.trim().replaceFirst(
      RegExp(r'/+$'),
      '',
    );
    final provider = _providerController.text.trim();
    if (baseUrl.isEmpty || provider.isEmpty) {
      setState(() {
        _error = 'Enter the API URL and provider name first.';
        _message = null;
      });
      return;
    }

    setState(() {
      _loadingSlots = true;
      _error = null;
      _message = null;
    });

    try {
      final uri = Uri.parse('$baseUrl/slots')
          .replace(queryParameters: {'provider': provider});
      final response = await _client
          .get(uri, headers: _headers)
          .timeout(const Duration(seconds: 15));
      final body = _decodeBody(response);
      if (response.statusCode < 200 || response.statusCode >= 300) {
        throw ApiException(_errorFromBody(body, response.statusCode));
      }

      final data = body as Map<String, dynamic>;
      final slots = (data['slots'] as List<dynamic>)
          .map((item) => Slot.fromJson(item as Map<String, dynamic>))
          .toList();
      setState(() {
        _slots = slots;
        _selectedSlot = null;
        _message = slots.isEmpty ? 'No slots found for $provider.' : null;
      });
    } on FormatException {
      setState(() => _error = 'Enter a valid API URL, including https://.');
    } on ApiException catch (error) {
      setState(() => _error = error.message);
    } catch (error) {
      setState(() => _error = 'Could not load slots: $error');
    } finally {
      if (mounted) setState(() => _loadingSlots = false);
    }
  }

  Future<void> _createBooking() async {
    final slot = _selectedSlot;
    final name = _nameController.text.trim();
    final email = _emailController.text.trim();
    final baseUrl = _apiUrlController.text.trim().replaceFirst(
      RegExp(r'/+$'),
      '',
    );
    if (slot == null) {
      setState(() => _error = 'Select an available time slot.');
      return;
    }
    if (name.isEmpty || !_isEmail(email)) {
      setState(() => _error = 'Enter your name and a valid email address.');
      return;
    }
    if (baseUrl.isEmpty) {
      setState(() => _error = 'Enter the API URL first.');
      return;
    }

    setState(() {
      _submitting = true;
      _error = null;
      _message = null;
    });

    try {
      final uri = Uri.parse('$baseUrl/bookings');
      final response = await _client
          .post(
            uri,
            headers: {..._headers, 'Content-Type': 'application/json'},
            body: jsonEncode({'id': slot.id, 'name': name, 'email': email}),
          )
          .timeout(const Duration(seconds: 15));
      final body = _decodeBody(response);
      if (response.statusCode < 200 || response.statusCode >= 300) {
        throw ApiException(_errorFromBody(body, response.statusCode));
      }

      final booking = body as Map<String, dynamic>;
      setState(() {
        _message =
            'Booking ${booking['status'] ?? 'created'} '
            '(${booking['id'] ?? ''}).';
        _selectedSlot = null;
        _nameController.clear();
        _emailController.clear();
      });
      await _loadSlots();
    } on FormatException {
      setState(() => _error = 'Enter a valid API URL, including https://.');
    } on ApiException catch (error) {
      setState(() => _error = error.message);
    } catch (error) {
      setState(() => _error = 'Could not create booking: $error');
    } finally {
      if (mounted) setState(() => _submitting = false);
    }
  }

  Future<void> _resetSlots() async {
    final baseUrl = _apiUrlController.text.trim().replaceFirst(
      RegExp(r'/+$'),
      '',
    );
    if (baseUrl.isEmpty) {
      setState(() => _error = 'Enter the API URL first.');
      return;
    }

    dynamic payload;
    try {
      payload = jsonDecode(_resetSlotsController.text);
      if (payload is! Map<String, dynamic> || payload['slots'] is! List) {
        throw const FormatException('Expected an object with a slots list.');
      }
      for (final item in payload['slots'] as List) {
        if (item is! Map<String, dynamic> ||
            !['id', 'provider', 'start', 'end'].every(
              (key) => item[key] is String && (item[key] as String).isNotEmpty,
            )) {
          throw const FormatException(
            'Each slot needs non-empty id, provider, start, and end fields.',
          );
        }
      }
    } on FormatException catch (error) {
      setState(() => _error = 'Invalid reset JSON: ${error.message}');
      return;
    }

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (dialogContext) => AlertDialog(
        icon: const Icon(Icons.warning_amber_rounded, color: Colors.red),
        title: const Text('Replace all slots?'),
        content: const Text(
          'This deletes every current slot and recreates only the slots in '
          'the JSON below. Existing bookings can also be deleted by the '
          'database cascade. This cannot be undone.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(dialogContext, false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            style: FilledButton.styleFrom(
              backgroundColor: const Color(0xFF9B2929),
            ),
            onPressed: () => Navigator.pop(dialogContext, true),
            child: const Text('Replace all slots'),
          ),
        ],
      ),
    );
    if (confirmed != true || !mounted) return;

    setState(() {
      _resettingSlots = true;
      _error = null;
      _message = null;
    });

    try {
      final response = await _client
          .post(
            Uri.parse('$baseUrl/reset'),
            headers: {..._headers, 'Content-Type': 'application/json'},
            body: jsonEncode(payload),
          )
          .timeout(const Duration(seconds: 15));
      final body = _decodeBody(response);
      if (response.statusCode < 200 || response.statusCode >= 300) {
        throw ApiException(_errorFromBody(body, response.statusCode));
      }

      setState(() {
        _slots = [];
        _selectedSlot = null;
        _message = 'Slots were replaced successfully.';
      });
    } on FormatException {
      setState(() => _error = 'Enter a valid API URL, including https://.');
    } on ApiException catch (error) {
      setState(() => _error = error.message);
    } catch (error) {
      setState(() => _error = 'Could not reset slots: $error');
    } finally {
      if (mounted) setState(() => _resettingSlots = false);
    }
  }

  Future<void> _loadBookings({String? cursor}) async {
    final baseUrl = _baseUrl;
    final email = _bookingEmailController.text.trim();
    if (baseUrl.isEmpty || !_isEmail(email)) {
      setState(() => _error = 'Enter the API URL and a valid booking email.');
      return;
    }

    setState(() {
      _loadingBookings = true;
      _error = null;
      _message = null;
    });

    try {
      final query = <String, String>{'email': email, 'limit': '10'};
      if (cursor != null) query['cursor'] = cursor;
      final response = await _client
          .get(
            Uri.parse('$baseUrl/bookings').replace(queryParameters: query),
            headers: _headers,
          )
          .timeout(const Duration(seconds: 15));
      final body = _decodeBody(response);
      if (response.statusCode < 200 || response.statusCode >= 300) {
        throw ApiException(_errorFromBody(body, response.statusCode));
      }

      final data = body as Map<String, dynamic>;
      final bookings = (data['bookings'] as List<dynamic>)
          .map((item) => Booking.fromJson(item as Map<String, dynamic>))
          .toList();
      setState(() {
        if (cursor == null) {
          _bookings = bookings;
        } else {
          _bookings = [..._bookings, ...bookings];
        }
        _nextBookingCursor = data['next_cursor'] as String?;
        _message = _bookings.isEmpty ? 'No bookings found for $email.' : null;
      });
    } on FormatException {
      setState(() => _error = 'The booking list returned invalid JSON.');
    } on ApiException catch (error) {
      setState(() => _error = error.message);
    } catch (error) {
      setState(() => _error = 'Could not load bookings: $error');
    } finally {
      if (mounted) setState(() => _loadingBookings = false);
    }
  }

  Future<void> _loadBookingById() async {
    final baseUrl = _baseUrl;
    final id = _bookingIdController.text.trim();
    if (baseUrl.isEmpty || id.isEmpty) {
      setState(() => _error = 'Enter the API URL and a booking ID.');
      return;
    }

    setState(() {
      _loadingBookingById = true;
      _error = null;
      _message = null;
    });

    try {
      final response = await _client
          .get(
            Uri.parse('$baseUrl/bookings/${Uri.encodeComponent(id)}'),
            headers: _headers,
          )
          .timeout(const Duration(seconds: 15));
      final body = _decodeBody(response);
      if (response.statusCode < 200 || response.statusCode >= 300) {
        throw ApiException(_errorFromBody(body, response.statusCode));
      }

      final booking = Booking.fromJson(body as Map<String, dynamic>);
      setState(() {
        _bookingIdController.text = booking.id;
        _bookings = [booking];
        _nextBookingCursor = null;
      });
    } on FormatException {
      setState(() => _error = 'The booking response was invalid.');
    } on ApiException catch (error) {
      setState(() => _error = error.message);
    } catch (error) {
      setState(() => _error = 'Could not load booking: $error');
    } finally {
      if (mounted) setState(() => _loadingBookingById = false);
    }
  }

  Future<void> _cancelBooking(Booking booking) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (dialogContext) => AlertDialog(
        title: const Text('Cancel this booking?'),
        content: Text(
          'Cancel booking ${booking.id} for ${booking.name}? '
          'The time slot will become available again.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(dialogContext, false),
            child: const Text('Keep booking'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(dialogContext, true),
            child: const Text('Cancel booking'),
          ),
        ],
      ),
    );
    if (confirmed != true || !mounted) return;

    setState(() {
      _cancellingBooking = true;
      _error = null;
      _message = null;
    });
    try {
      final response = await _client
          .delete(
            Uri.parse('$_baseUrl/bookings/${Uri.encodeComponent(booking.id)}'),
            headers: _headers,
          )
          .timeout(const Duration(seconds: 15));
      final body = _decodeBody(response);
      if (response.statusCode < 200 || response.statusCode >= 300) {
        throw ApiException(_errorFromBody(body, response.statusCode));
      }

      final updatedBooking = Booking.fromJson(body as Map<String, dynamic>);
      setState(() {
        _bookings = _bookings
            .map((item) => item.id == updatedBooking.id ? updatedBooking : item)
            .toList();
        _message = 'Booking ${updatedBooking.id} was cancelled.';
      });
    } on FormatException {
      setState(() => _error = 'The cancellation response was invalid.');
    } on ApiException catch (error) {
      setState(() => _error = error.message);
    } catch (error) {
      setState(() => _error = 'Could not cancel booking: $error');
    } finally {
      if (mounted) setState(() => _cancellingBooking = false);
    }
  }

  Map<String, String> get _headers => {
    if (_apiKeyController.text.trim().isNotEmpty)
      'X-API-Key': _apiKeyController.text.trim(),
  };

  String get _baseUrl =>
      _apiUrlController.text.trim().replaceFirst(RegExp(r'/+$'), '');

  bool _isEmail(String value) =>
      RegExp(r'^[^@\s]+@[^@\s]+\.[^@\s]+$').hasMatch(value);

  dynamic _decodeBody(http.Response response) {
    if (response.body.isEmpty) return <String, dynamic>{};
    try {
      return jsonDecode(response.body);
    } on FormatException {
      return response.body;
    }
  }

  String _errorFromBody(dynamic body, int statusCode) {
    if (body is Map<String, dynamic> && body['error'] is String) {
      return body['error'] as String;
    }
    return 'The API returned HTTP $statusCode.';
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Row(
          children: [
            Icon(Icons.event_available, color: Color(0xFF167D73)),
            SizedBox(width: 10),
            Text('Booking'),
          ],
        ),
      ),
      body: LayoutBuilder(
        builder: (context, constraints) {
          final wide = constraints.maxWidth >= 760;
          final content = wide
              ? Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Expanded(flex: 5, child: _buildSlotPanel()),
                    const SizedBox(width: 24),
                    Expanded(flex: 4, child: _buildBookingPanel()),
                  ],
                )
              : Column(
                  children: [
                    _buildSlotPanel(),
                    const SizedBox(height: 20),
                    _buildBookingPanel(),
                  ],
                );

          return Center(
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 1120),
              child: ListView(
                padding: EdgeInsets.symmetric(
                  horizontal: wide ? 32 : 18,
                  vertical: 18,
                ),
                children: [
                  _buildIntro(),
                  const SizedBox(height: 22),
                  _buildConnectionPanel(),
                  const SizedBox(height: 22),
                  if (_error != null) _buildNotice(_error!, isError: true),
                  if (_message != null) _buildNotice(_message!),
                  if (_error != null || _message != null)
                    const SizedBox(height: 16),
                  content,
                  const SizedBox(height: 22),
                  _buildBookingsPanel(),
                  const SizedBox(height: 22),
                  _buildResetPanel(),
                ],
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildIntro() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Find a time that works.',
          style: Theme.of(context).textTheme.headlineMedium?.copyWith(
            fontWeight: FontWeight.w700,
            color: const Color(0xFF182B2A),
          ),
        ),
        const SizedBox(height: 6),
        Text(
          'Connect to the booking API, choose a provider, and reserve an open slot.',
          style: Theme.of(context).textTheme.bodyLarge
              ?.copyWith(color: const Color(0xFF5A6A65)),
        ),
      ],
    );
  }

  Widget _buildConnectionPanel() {
    return _Panel(
      title: 'API connection',
      subtitle:
          'The API key is used only in this browser session. Do not put a '
          'private production key in a publicly hosted frontend.',
      child: LayoutBuilder(
        builder: (context, constraints) {
          final vertical = constraints.maxWidth < 620;
          final fields = [
            _textField(
              controller: _apiUrlController,
              label: 'API base URL',
              hint: 'https://your-api.onrender.com',
              keyboardType: TextInputType.url,
            ),
            _textField(
              controller: _apiKeyController,
              label: 'API key',
              hint: 'X-API-Key',
              obscureText: true,
            ),
            _textField(
              controller: _providerController,
              label: 'Provider',
              hint: 'Provider name',
              onSubmitted: (_) => _loadSlots(),
            ),
          ];

          if (vertical) {
            return Column(
              children: [
                ...fields.map(
                  (field) => Padding(
                    padding: const EdgeInsets.only(bottom: 12),
                    child: field,
                  ),
                ),
                Align(alignment: Alignment.centerRight, child: _loadButton()),
              ],
            );
          }

          return Column(
            children: [
              Row(
                children: [
                  Expanded(child: fields[0]),
                  const SizedBox(width: 12),
                  Expanded(child: fields[1]),
                ],
              ),
              const SizedBox(height: 12),
              Row(
                children: [
                  Expanded(child: fields[2]),
                  const SizedBox(width: 12),
                  _loadButton(),
                ],
              ),
            ],
          );
        },
      ),
    );
  }

  Widget _loadButton() => FilledButton.icon(
    onPressed: _loadingSlots ? null : _loadSlots,
    icon: _loadingSlots
        ? const SizedBox(
            width: 18,
            height: 18,
            child: CircularProgressIndicator(strokeWidth: 2),
          )
        : const Icon(Icons.search),
    label: Text(_loadingSlots ? 'Loading' : 'Find slots'),
  );

  Widget _buildSlotPanel() {
    return _Panel(
      title: 'Available slots',
      subtitle: _slots.isEmpty
          ? 'Search by provider to see open appointment times.'
          : '${_slots.where((slot) => slot.available).length} open time(s)',
      child: _slots.isEmpty
          ? const _EmptyState(
              icon: Icons.calendar_month_outlined,
              text: 'Your available times will appear here.',
            )
          : ListView.separated(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              itemCount: _slots.length,
              separatorBuilder: (context, index) => const SizedBox(height: 10),
              itemBuilder: (context, index) {
                final slot = _slots[index];
                final selected = _selectedSlot?.id == slot.id;
                return _SlotTile(
                  slot: slot,
                  selected: selected,
                  onTap: slot.available
                      ? () => setState(() {
                          _selectedSlot = slot;
                          _error = null;
                        })
                      : null,
                );
              },
            ),
    );
  }

  Widget _buildBookingPanel() {
    return _Panel(
      title: 'Your details',
      subtitle: _selectedSlot == null
          ? 'Select an available time before booking.'
          : 'Selected: ${_selectedSlot!.provider} · ${_formatTime(_selectedSlot!.start)}',
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          _textField(
            controller: _nameController,
            label: 'Full name',
            hint: 'Your name',
            textCapitalization: TextCapitalization.words,
          ),
          const SizedBox(height: 12),
          _textField(
            controller: _emailController,
            label: 'Email address',
            hint: 'you@example.com',
            keyboardType: TextInputType.emailAddress,
          ),
          const SizedBox(height: 18),
          FilledButton.icon(
            onPressed: _submitting ? null : _createBooking,
            icon: _submitting
                ? const SizedBox(
                    width: 18,
                    height: 18,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Icon(Icons.check_circle_outline),
            label: Text(_submitting ? 'Booking…' : 'Confirm booking'),
            style: FilledButton.styleFrom(
              padding: const EdgeInsets.symmetric(vertical: 16),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildResetPanel() {
    return _Panel(
      title: 'Administrative: replace all slots',
      subtitle:
          'POST /reset replaces the complete slot schedule. Existing bookings '
          'may be deleted by the database cascade.',
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          TextField(
            key: const Key('reset-slots-json'),
            controller: _resetSlotsController,
            minLines: 8,
            maxLines: 16,
            keyboardType: TextInputType.multiline,
            decoration: const InputDecoration(
              labelText: 'Replacement slots JSON',
              alignLabelWithHint: true,
              hintText: '{"slots": [{"id": "...", "provider": "...", "start": "...", "end": "..."}]}',
            ),
            style: const TextStyle(fontFamily: 'monospace', fontSize: 13),
          ),
          const SizedBox(height: 12),
          Align(
            alignment: Alignment.centerRight,
            child: OutlinedButton.icon(
              key: const Key('replace-all-slots'),
              onPressed: _resettingSlots ? null : _resetSlots,
              icon: _resettingSlots
                  ? const SizedBox(
                      width: 18,
                      height: 18,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : const Icon(Icons.delete_sweep_outlined),
              label: Text(_resettingSlots ? 'Replacing…' : 'Replace all slots'),
              style: OutlinedButton.styleFrom(
                foregroundColor: const Color(0xFF9B2929),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBookingsPanel() {
    return _Panel(
      title: 'Bookings',
      subtitle: 'List bookings by email, look up a booking by ID, or cancel a booking.',
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          LayoutBuilder(
            builder: (context, constraints) {
              final vertical = constraints.maxWidth < 620;
              final emailField = _textField(
                key: const Key('booking-email-search'),
                controller: _bookingEmailController,
                label: 'Booking email',
                hint: 'you@example.com',
                keyboardType: TextInputType.emailAddress,
                onSubmitted: (_) => _loadBookings(),
              );
              final listButton = FilledButton.tonalIcon(
                key: const Key('list-bookings'),
                onPressed: _loadingBookings ? null : () => _loadBookings(),
                icon: _loadingBookings
                    ? const SizedBox(
                        width: 18,
                        height: 18,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Icon(Icons.list_alt),
                label: const Text('List bookings'),
              );
              final idField = _textField(
                key: const Key('booking-id-search'),
                controller: _bookingIdController,
                label: 'Booking ID',
                hint: 'UUID',
                onSubmitted: (_) => _loadBookingById(),
              );
              final lookupButton = OutlinedButton.icon(
                key: const Key('find-booking-by-id'),
                onPressed: _loadingBookingById ? null : _loadBookingById,
                icon: _loadingBookingById
                    ? const SizedBox(
                        width: 18,
                        height: 18,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Icon(Icons.search),
                label: const Text('Find by ID'),
              );

              if (vertical) {
                return Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    emailField,
                    const SizedBox(height: 10),
                    Align(alignment: Alignment.centerRight, child: listButton),
                    const SizedBox(height: 16),
                    idField,
                    const SizedBox(height: 10),
                    Align(
                      alignment: Alignment.centerRight,
                      child: lookupButton,
                    ),
                  ],
                );
              }
              return Column(
                children: [
                  Row(
                    children: [
                      Expanded(child: emailField),
                      const SizedBox(width: 12),
                      listButton,
                    ],
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(child: idField),
                      const SizedBox(width: 12),
                      lookupButton,
                    ],
                  ),
                ],
              );
            },
          ),
          const SizedBox(height: 16),
          if (_bookings.isEmpty)
            const _EmptyState(
              icon: Icons.event_note_outlined,
              text: 'Search by email or booking ID to see bookings.',
            )
          else
            ListView.separated(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              itemCount: _bookings.length,
              separatorBuilder: (context, index) => const SizedBox(height: 10),
              itemBuilder: (context, index) {
                final booking = _bookings[index];
                return _BookingTile(
                  booking: booking,
                  cancelling: _cancellingBooking,
                  onCancel: booking.status.toLowerCase() == 'cancelled'
                      ? null
                      : () => _cancelBooking(booking),
                );
              },
            ),
          if (_nextBookingCursor != null) ...[
            const SizedBox(height: 12),
            Align(
              alignment: Alignment.center,
              child: OutlinedButton.icon(
                onPressed: _loadingBookings
                    ? null
                    : () => _loadBookings(cursor: _nextBookingCursor),
                icon: const Icon(Icons.expand_more),
                label: const Text('Load more'),
              ),
            ),
          ],
        ],
      ),
    );
  }

  Widget _textField({
    Key? key,
    required TextEditingController controller,
    required String label,
    required String hint,
    TextInputType? keyboardType,
    TextCapitalization textCapitalization = TextCapitalization.none,
    bool obscureText = false,
    ValueChanged<String>? onSubmitted,
  }) {
    return TextField(
      key: key,
      controller: controller,
      decoration: InputDecoration(labelText: label, hintText: hint),
      keyboardType: keyboardType,
      textCapitalization: textCapitalization,
      obscureText: obscureText,
      onSubmitted: onSubmitted,
    );
  }

  Widget _buildNotice(String message, {bool isError = false}) {
    final color = isError ? const Color(0xFF9B2929) : const Color(0xFF176B55);
    return Container(
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: isError ? const Color(0xFFFFEEEE) : const Color(0xFFE9F6EF),
        borderRadius: BorderRadius.circular(14),
      ),
      child: Row(
        children: [
          Icon(
            isError ? Icons.error_outline : Icons.check_circle_outline,
            color: color,
          ),
          const SizedBox(width: 10),
          Expanded(
            child: Text(message, style: TextStyle(color: color)),
          ),
        ],
      ),
    );
  }

  String _formatTime(String value) {
    final parsed = DateTime.tryParse(value)?.toLocal();
    if (parsed == null) return value;
    final date = '${_month(parsed.month)} ${parsed.day}';
    final hour = parsed.hour % 12 == 0 ? 12 : parsed.hour % 12;
    final minute = parsed.minute.toString().padLeft(2, '0');
    final period = parsed.hour < 12 ? 'AM' : 'PM';
    return '$date · $hour:$minute $period';
  }

  String _month(int month) => const [
    'Jan',
    'Feb',
    'Mar',
    'Apr',
    'May',
    'Jun',
    'Jul',
    'Aug',
    'Sep',
    'Oct',
    'Nov',
    'Dec',
  ][month - 1];
}

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

class ApiException implements Exception {
  const ApiException(this.message);
  final String message;
}

class _BookingTile extends StatelessWidget {
  const _BookingTile({
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

class _Panel extends StatelessWidget {
  const _Panel({
    required this.title,
    required this.subtitle,
    required this.child,
  });

  final String title;
  final String subtitle;
  final Widget child;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: const Color(0xFFE4E9E3)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            title,
            style: Theme.of(context).textTheme.titleLarge
                ?.copyWith(fontWeight: FontWeight.w700),
          ),
          const SizedBox(height: 4),
          Text(
            subtitle,
            style: Theme.of(context).textTheme.bodyMedium
                ?.copyWith(color: const Color(0xFF64736D)),
          ),
          const SizedBox(height: 18),
          child,
        ],
      ),
    );
  }
}

class _SlotTile extends StatelessWidget {
  const _SlotTile({
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

class _EmptyState extends StatelessWidget {
  const _EmptyState({required this.icon, required this.text});

  final IconData icon;
  final String text;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.symmetric(vertical: 30, horizontal: 16),
      decoration: BoxDecoration(
        color: const Color(0xFFFAFBF8),
        borderRadius: BorderRadius.circular(14),
      ),
      child: Column(
        children: [
          Icon(icon, size: 30, color: const Color(0xFF7B8B84)),
          const SizedBox(height: 10),
          Text(text, textAlign: TextAlign.center),
        ],
      ),
    );
  }
}
