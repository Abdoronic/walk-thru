package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;

public class CustomerViewShopsActivity extends AppCompatActivity {


    private Button ordersButton;
    private Button fakeItemButton;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_customer_view_shops);

        ordersButton = findViewById(R.id.ordersButton);
        fakeItemButton = findViewById(R.id.fakeItemButton);

        ordersButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent i = new Intent(CustomerViewShopsActivity.this, CustomerViewOrdersActivity.class);
                i.putExtra("activity","shops");
                startActivity(i);
            }
        });
        fakeItemButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent i = new Intent(CustomerViewShopsActivity.this, CustomerViewItemsActivity.class);
                startActivity(i);
            }
        });

    }
}
