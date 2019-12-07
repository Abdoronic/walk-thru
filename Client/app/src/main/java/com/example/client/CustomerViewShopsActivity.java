package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.ListView;
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

import java.util.HashMap;
import java.util.Map;

public class CustomerViewShopsActivity extends AppCompatActivity {


    private Button ordersButton;

    private String[] shopData;
    private JSONObject[] shopDataJSON;
    private ArrayAdapter<String> myAdapter;
    private ListView shopsListView;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_customer_view_shops);

        ordersButton = findViewById(R.id.ordersButton);
        shopsListView = findViewById(R.id.shopsListView);

        ordersButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent i = new Intent(CustomerViewShopsActivity.this, CustomerViewOrdersActivity.class);
                i.putExtra("activity","shops");
                i.putExtra("id",getIntent().getIntExtra("id",-1));
                startActivity(i);
            }
        });


        RequestQueue queue = Volley.newRequestQueue(CustomerViewShopsActivity.this);

        String url = "http://10.0.2.2:8000/shops";
        JsonArrayRequest jsonArrayRequest = new JsonArrayRequest(Request.Method.GET, url, null, new Response.Listener<JSONArray>() {
            @Override
            public void onResponse(JSONArray response) {
                try {
                    shopData = new String[response.length()];
                    shopDataJSON = new JSONObject[response.length()];
                    for(int i=0;i<response.length();i++){
                        JSONObject shopJSON = response.getJSONObject(i);
                        String shop = "Shop Name: "+shopJSON.getString("name")+"\n"+"Shop Location: "+shopJSON.getString("location");
                        shopData[i]=shop;
                        shopDataJSON[i]=shopJSON;
                    }
                    myAdapter = new ArrayAdapter<String>(CustomerViewShopsActivity.this, android.R.layout.simple_list_item_1, android.R.id.text1, shopData);
                    shopsListView.setAdapter(myAdapter);
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



        shopsListView.setOnItemClickListener(new AdapterView.OnItemClickListener() {

            @Override
            public void onItemClick(AdapterView<?> parent, View view, final int position, long id) {
                RequestQueue queue = Volley.newRequestQueue(CustomerViewShopsActivity.this);
                JSONObject jsonBody = new JSONObject();
                try {
                    jsonBody.put("delivered",false);
                } catch (JSONException e) {
                    e.printStackTrace();
                }
                String url = "http://10.0.2.2:8000/customers/"+getIntent().getIntExtra("id",-1)+"/createOrder";
                JsonObjectRequest jsonObjectRequest = new JsonObjectRequest(Request.Method.POST, url, jsonBody, new Response.Listener<JSONObject>() {
                    @Override
                    public void onResponse(JSONObject response) {
                        try {
                            int shopID = shopDataJSON[position].getInt("id");
                            Intent i = new Intent(getApplicationContext(),CustomerViewItemsActivity.class);
                            i.putExtra("customerID", getIntent().getIntExtra("id",-1));
                            i.putExtra("shopID",shopID);
                            i.putExtra("orderID",response.getInt("id"));
                            startActivity(i);
                        } catch (JSONException e) {
                            e.printStackTrace();
                        }
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
