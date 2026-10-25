import 'package:booking_web/page/booking_page.dart';
import 'package:flutter/material.dart';

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
