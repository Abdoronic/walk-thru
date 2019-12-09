package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.ListView;
import android.widget.Toast;

import com.android.volley.Request;
import com.android.volley.RequestQueue;
import com.android.volley.Response;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.JsonArrayRequest;
import com.android.volley.toolbox.Volley;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

public class CustomerViewOrdersActivity extends AppCompatActivity {

    private Button shopButton;
    private String[] orderData;
    private JSONObject[] orderDataJSON;
    private ArrayAdapter<String> myAdapter;
    private ListView ordersListView;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_customer_view_orders);

        shopButton = findViewById(R.id.shopsButton);
        ordersListView = findViewById(R.id.ordersListView);

        if(getIntent().getStringExtra("activity").equals("items")){
            shopButton.setText("Items");
            shopButton.setOnClickListener(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    Intent i = new Intent(CustomerViewOrdersActivity.this, CustomerViewItemsActivity.class);
                    i.putExtra("id", getIntent().getIntExtra("id",-1));
                    i.putExtra("firstName", getIntent().getStringExtra("firstName"));
                    i.putExtra("lastName", getIntent().getStringExtra("lastName"));
                    i.putExtra("email", getIntent().getStringExtra("email"));
                    i.putExtra("password", getIntent().getStringExtra("password"));
                    i.putExtra("creditCardNumber", getIntent().getStringExtra("creditCardNumber"));
                    i.putExtra("creditCardExpiryDate", getIntent().getStringExtra("creditCardExpiryDate"));
                    i.putExtra("creditCardCVV", getIntent().getIntExtra("creditCardCVV",-1));

                    i.putExtra("shopID",getIntent().getIntExtra("shopID",-1));
                    i.putExtra("orderID",getIntent().getIntExtra("orderID",-1));
                    startActivity(i);
                }
            });
        }
        else{
            shopButton.setText("Shops");
            shopButton.setOnClickListener(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    Intent i = new Intent(CustomerViewOrdersActivity.this, CustomerViewShopsActivity.class);
                    i.putExtra("id", getIntent().getIntExtra("id",-1));
                    i.putExtra("firstName", getIntent().getStringExtra("firstName"));
                    i.putExtra("lastName", getIntent().getStringExtra("lastName"));
                    i.putExtra("email", getIntent().getStringExtra("email"));
                    i.putExtra("password", getIntent().getStringExtra("password"));
                    i.putExtra("creditCardNumber", getIntent().getStringExtra("creditCardNumber"));
                    i.putExtra("creditCardExpiryDate", getIntent().getStringExtra("creditCardExpiryDate"));
                    i.putExtra("creditCardCVV", getIntent().getIntExtra("creditCardCVV",-1));
                    startActivity(i);
                }
            });
        }

        RequestQueue queue = Volley.newRequestQueue(CustomerViewOrdersActivity.this);

        String url = "http://10.0.2.2:8000/customers/"+getIntent().getIntExtra("id",-1)+"/viewOrders";
        JsonArrayRequest jsonArrayRequest = new JsonArrayRequest(Request.Method.GET, url, null, new Response.Listener<JSONArray>() {
            @Override
            public void onResponse(JSONArray response) {
                try {
                    orderData = new String[response.length()];
                    orderDataJSON = new JSONObject[response.length()];
                    for(int i=0;i<response.length();i++){
                        JSONObject orderJSON = response.getJSONObject(i);
                        String order = "Order Number: " + orderJSON.getString("id") + "\n"
                                + "Delivered: " + orderJSON.getString("delivered") + "\n"
                                +  "Price: " +  orderJSON.getString("price") + "\n"
                                + "Date: " + orderJSON.getString("date") + '\n';
                        orderData[i]=order;
                        orderDataJSON[i]=orderJSON;
                    }
                    myAdapter = new ArrayAdapter<String>(CustomerViewOrdersActivity.this, android.R.layout.simple_list_item_1, android.R.id.text1, orderData);
                    ordersListView.setAdapter(myAdapter);
                } catch (JSONException e) {
                    e.printStackTrace();
                }
            }
        }, new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                try {
                    JSONObject errData =new JSONObject(new String(error.networkResponse.data));
                    Toast.makeText(getApplicationContext(),errData.getString("error"),Toast.LENGTH_LONG).show();

                } catch (JSONException e) {
                    e.printStackTrace();
                }
                error.printStackTrace();
            }
        });
        queue.add(jsonArrayRequest);


    }
}
