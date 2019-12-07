package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;

public class CustomerViewItemsActivity extends AppCompatActivity {

    private Button orderButton;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_customer_view_items);

        orderButton = findViewById(R.id.ordersButton);

        orderButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent i = new Intent(CustomerViewItemsActivity.this, CustomerViewOrdersActivity.class);
                i.putExtra("activity","items");
                startActivity(i);
            }
        });
    }
}
