package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.BaseAdapter;
import android.widget.Button;
import android.widget.ListAdapter;
import android.widget.ListView;
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

import java.util.ArrayList;


public class CustomerViewItemsActivity extends AppCompatActivity {

    private Button orderButton;
    private Button checkoutButton;

    private String[] itemData;
    private ArrayList<String> itemDataList;
    private JSONObject[] itemDataJSON;
//    private ArrayAdapter<String> myAdapter;
    private MyCustomAdapter adapter;
    private ListView itemsListView;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_customer_view_items);

        orderButton = findViewById(R.id.ordersButton);
        checkoutButton = findViewById(R.id.checkoutButton);
        itemsListView = findViewById(R.id.itemsListView);
        itemDataList = new ArrayList<String>();


        orderButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent i = new Intent(CustomerViewItemsActivity.this, CustomerViewOrdersActivity.class);
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

                i.putExtra("activity","items");
                startActivity(i);
            }
        });

        RequestQueue queue = Volley.newRequestQueue(CustomerViewItemsActivity.this);
        String url = "http://10.0.2.2:8000/customers/viewItems/"+getIntent().getIntExtra("shopID",-1);
        JsonArrayRequest jsonArrayRequest = new JsonArrayRequest(Request.Method.GET, url, null, new Response.Listener<JSONArray>() {
            @Override
            public void onResponse(JSONArray response) {
                try {
                    itemData = new String[response.length()];
                    itemDataJSON = new JSONObject[response.length()];
                    for(int i=0;i<response.length();i++){
                        JSONObject itemJSON = response.getJSONObject(i);
                        String item = "Item Name: "+itemJSON.getString("name")+"\n"+"Item Type: "+itemJSON.getString("type")+"\n"+"Item Description: "+itemJSON.getString("description")+"\n"+"Item Price: "+itemJSON.getDouble("price")+"\n";
                        itemData[i]=item;
                        itemDataJSON[i]=itemJSON;
                        itemDataList.add(item);
                    }
                    adapter = new MyCustomAdapter(itemDataList, CustomerViewItemsActivity.this);
                    itemsListView.setAdapter(adapter);

//                    myAdapter = new ArrayAdapter<String>(CustomerViewItemsActivity.this, android.R.layout.simple_list_item_1, android.R.id.text1, itemData);
//                    itemsListView.setAdapter(myAdapter);
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

        checkoutButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                int []itemsQuantities = adapter.getQuantity();
                for(int i=0;i<itemsQuantities.length;i++){
                    if(itemsQuantities[i]>0){
                        RequestQueue queue = Volley.newRequestQueue(CustomerViewItemsActivity.this);
                        String url = null;
                        try {
                            url = "http://10.0.2.2:8000/customers/"+getIntent().getIntExtra("id",-1)+"/addItem/"+getIntent().getIntExtra("orderID",-1)+"/"+itemDataJSON[i].getInt("id")+"/"+itemsQuantities[i];
                        } catch (JSONException e) {
                            e.printStackTrace();
                        }
                        JsonObjectRequest jsonObjectRequest = new JsonObjectRequest(Request.Method.PUT, url, null, new Response.Listener<JSONObject>() {
                            @Override
                            public void onResponse(JSONObject response) {

                            }
                        }, new Response.ErrorListener() {
                            @Override
                            public void onErrorResponse(VolleyError error) {
                                error.printStackTrace();
                            }
                        });
                        queue.add(jsonObjectRequest);
                    }
                }
                Intent i = new Intent(CustomerViewItemsActivity.this, CustomerOrderSummaryActivity.class);
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

//        itemsListView.setOnItemClickListener(new AdapterView.OnItemClickListener() {
//
//            @Override
//            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
//                try {
//                    int itemID = itemDataJSON[position].getInt("id");
//                    int shopID = itemDataJSON[position].getInt("shopID");
//                    Intent i = new Intent(getApplicationContext(),CustomerItemClickActivity.class);
//                    i.putExtra("itemID",itemID);
//                    i.putExtra("shopID",shopID);
//                    i.putExtra("customerID",getIntent().getIntExtra("customerID",-1));
//                    i.putExtra("orderID",getIntent().getIntExtra("orderID",-1));
//                    startActivity(i);
//                } catch (JSONException e) {
//                    e.printStackTrace();
//                }
//            }
//        });
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        RequestQueue queue = Volley.newRequestQueue(CustomerViewItemsActivity.this);
        String url = "http://10.0.2.2:8000/orders/"+getIntent().getIntExtra("orderID",-1);
        JsonObjectRequest jsonObjectRequest = new JsonObjectRequest(Request.Method.GET, url, null, new Response.Listener<JSONObject>() {
            @Override
            public void onResponse(JSONObject response) {
                if(response.isNull("shopID")) {
                    RequestQueue queue = Volley.newRequestQueue(CustomerViewItemsActivity.this);
                    String url = "http://10.0.2.2:8000/orders/" + getIntent().getIntExtra("orderID", -1);
                    JsonObjectRequest jsonObjectRequest = new JsonObjectRequest(Request.Method.DELETE, url, null, new Response.Listener<JSONObject>() {
                        @Override
                        public void onResponse(JSONObject response) {
                        }
                    },
                            new Response.ErrorListener() {
                                @Override
                                public void onErrorResponse(VolleyError error) {
                                    error.printStackTrace();
                                }
                            });
                    queue.add(jsonObjectRequest);
                }
            }
        },
        new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                error.printStackTrace();
            }
        });
        queue.add(jsonObjectRequest);
    }
}
