package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.os.Bundle;
import android.widget.TextView;

public class CustomerViewShops extends AppCompatActivity {

    private TextView welcomeCustomerTextView;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_customer_view_shops);
        welcomeCustomerTextView = findViewById(R.id.welcomeCustomerTextView);
        welcomeCustomerTextView.setText("Welcome, "+getIntent().getStringExtra("firstName"));
    }
}
