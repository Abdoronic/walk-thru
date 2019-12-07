package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.TextClock;
import android.widget.TextView;
import android.widget.Toast;

import com.android.volley.Request;
import com.android.volley.RequestQueue;
import com.android.volley.Response;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.JsonArrayRequest;
import com.android.volley.toolbox.JsonObjectRequest;
import com.android.volley.toolbox.Volley;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

public class CustomerItemClickActivity extends AppCompatActivity {

    private Button addButton;
    private Button removeButton;
    private TextView quantityTextView;
    private int quantity = 0;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_customer_item_click);

        addButton = findViewById(R.id.addButton);
        removeButton = findViewById(R.id.removeButton);
        quantityTextView = findViewById(R.id.quantityTextView);

        RequestQueue queue = Volley.newRequestQueue(CustomerItemClickActivity.this);
        String url = "http://10.0.2.2:8000/customers/" + getIntent().getIntExtra("customerID", -1) + "/viewOrderItems/" + getIntent().getIntExtra("orderID",-1);
        JsonArrayRequest jsonArrayRequest = new JsonArrayRequest(Request.Method.GET, url, null, new Response.Listener<JSONArray>() {
            @Override
            public void onResponse(JSONArray response) {
                for(int i = 0 ;i<response.length();i++){
                    try {
                        JSONObject item = response.getJSONObject(i);
                        if(item.getInt("itemID") == getIntent().getIntExtra("itemID",-1)){
                            quantity = item.getInt("quantity");
                            break;
                        }
                    } catch (JSONException e) {
                        e.printStackTrace();
                    }
                    quantityTextView.setText(quantity+"");
                }
            }
        }, new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                quantityTextView.setText("0");
                error.printStackTrace();
            }
        });
        queue.add(jsonArrayRequest);

        addButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                RequestQueue queue = Volley.newRequestQueue(CustomerItemClickActivity.this);
                String url = "http://10.0.2.2:8000/customers/" + getIntent().getIntExtra("customerID", -1) + "/addItem/" + getIntent().getIntExtra("orderID",-1) + "/" + getIntent().getIntExtra("itemID", -1);
                JsonObjectRequest jsonObjectRequest = new JsonObjectRequest(Request.Method.PUT, url, null, new Response.Listener<JSONObject>() {
                    @Override
                    public void onResponse(JSONObject response) {
                        int currentQuantity = Integer.parseInt(quantityTextView.getText().toString());
                        currentQuantity++;
                        quantityTextView.setText(currentQuantity+"");
                    }
                }, new Response.ErrorListener() {
                    @Override
                    public void onErrorResponse(VolleyError error) {
//                        try {
//                            JSONObject errData = new JSONObject(new String(error.networkResponse.data));
//                            Toast.makeText(getApplicationContext(), errData.getString("error"), Toast.LENGTH_LONG).show();
//
//                        } catch (JSONException e) {
//                            e.printStackTrace();
//                        }
                        error.printStackTrace();
                    }
                });
                queue.add(jsonObjectRequest);
            }
        });

        if(quantity>0) {
            removeButton.setOnClickListener(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    RequestQueue queue = Volley.newRequestQueue(CustomerItemClickActivity.this);

                    String url = "http://10.0.2.2:8000/customers/" + getIntent().getIntExtra("customerID", -1) + "/removeItem/" + getIntent().getIntExtra("orderID", -1) + "/" + getIntent().getIntExtra("itemID", -1);
                    JsonObjectRequest jsonObjectRequest = new JsonObjectRequest(Request.Method.PUT, url, null, new Response.Listener<JSONObject>() {
                        @Override
                        public void onResponse(JSONObject response) {
                            int currentQuantity = Integer.parseInt(quantityTextView.getText().toString());
                            currentQuantity--;
                            quantityTextView.setText(currentQuantity+"");
                        }
                    }, new Response.ErrorListener() {
                        @Override
                        public void onErrorResponse(VolleyError error) {
                            try {
                                JSONObject errData = new JSONObject(new String(error.networkResponse.data));
                                Toast.makeText(getApplicationContext(), errData.getString("error"), Toast.LENGTH_LONG).show();

                            } catch (JSONException e) {
                                e.printStackTrace();
                            }
                            error.printStackTrace();
                        }
                    });
                    queue.add(jsonObjectRequest);
                }
            });
        }
    }
}
