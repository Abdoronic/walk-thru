package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.ListView;
import android.widget.TextView;
import android.widget.Toast;

import com.android.volley.DefaultRetryPolicy;
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

public class CustomerOrderSummaryActivity extends AppCompatActivity {
    private static final String TAG = "CustomerOrderSummaryAct";
    private ArrayAdapter<String> summaryAdapter;
    private String[] orderData;
    private ListView summaryListView;
    private TextView customerNameTextView;
    private TextView orderIDTextView;
    private TextView amountTextView;
    private Button confirmButton;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_customer_order_summary);

        summaryListView = findViewById(R.id.orderListView);
        customerNameTextView= findViewById(R.id.customerNameTextView);
        orderIDTextView = findViewById(R.id.orderIDTextView);
        amountTextView= findViewById(R.id.amountTextView);
        customerNameTextView.setText(getIntent().getStringExtra("firstName"));
        orderIDTextView.setText(getIntent().getIntExtra("orderID",-1)+"");
        confirmButton = findViewById(R.id.confirmButton);

        RequestQueue queue = Volley.newRequestQueue(CustomerOrderSummaryActivity.this);
        String url = "http://10.0.2.2:8000/customers/"+getIntent().getIntExtra("id",-1)+"/viewOrderItems"+"/"+getIntent().getIntExtra("orderID",-1);
        JsonArrayRequest jsonArrayRequest = new JsonArrayRequest(Request.Method.GET, url, null, new Response.Listener<JSONArray>() {
            @Override
            public void onResponse(JSONArray response) {
                try {
                    orderData = new String[response.length()];
                    for (int i = 0; i < response.length(); i++) {
                        JSONObject orderItemJSON = response.getJSONObject(i);
                        String orderItem = "x"+orderItemJSON.getInt("quantity") + "  "+ orderItemJSON.getString("name")+"  "+orderItemJSON.getString("type")+"  "+orderItemJSON.getDouble("price");
                        orderData[i]=orderItem;
                        amountTextView.setText(orderItemJSON.getDouble("OrderPrice")+"");
                        Log.d(TAG, "onResponse: "+orderItemJSON.getDouble("OrderPrice"));
                    }
                    summaryAdapter = new ArrayAdapter<String>(CustomerOrderSummaryActivity.this, android.R.layout.simple_list_item_1, android.R.id.text1, orderData);
                    summaryListView.setAdapter(summaryAdapter);
                }catch (JSONException e) {
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

        confirmButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                RequestQueue queue = Volley.newRequestQueue(CustomerOrderSummaryActivity.this);
                String url = "http://10.0.2.2:8000/customers/"+getIntent().getIntExtra("id",-1)+"/checkout"+"/"+getIntent().getIntExtra("orderID",-1)+"/"+getIntent().getIntExtra("shopID",-1);
                JsonObjectRequest jsonObjectRequest = new JsonObjectRequest(Request.Method.POST, url, null, new Response.Listener<JSONObject>() {
                    @Override
                    public void onResponse(JSONObject response) {
                        Toast.makeText(getApplicationContext(),"Success!",Toast.LENGTH_LONG).show();
//                        Intent i = new Intent(getApplicationContext(), CustomerViewShopsActivity.class);
//                        //i.setFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP);
//                        i.setFlags(Intent.FLAG_ACTIVITY_CLEAR_TASK | Intent.FLAG_ACTIVITY_NEW_TASK | Intent.FLAG_ACTIVITY_CLEAR_TOP);
//                        i.putExtra("id",  getIntent().getIntExtra("id",-1));
//                        i.putExtra("firstName", getIntent().getStringExtra("firstName"));
//                        i.putExtra("lastName", getIntent().getStringExtra("lastName"));
//                        i.putExtra("email", getIntent().getStringExtra("email"));
//                        i.putExtra("password", getIntent().getStringExtra("password"));
//                        i.putExtra("creditCardNumber", getIntent().getStringExtra("creditCardNumber"));
//                        i.putExtra("creditCardExpiryDate", getIntent().getStringExtra("creditCardExpiryDate"));
//                        i.putExtra("creditCardCVV", getIntent().getIntExtra("creditCardCVV",-1));
//                        startActivity(i);
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
                jsonObjectRequest.setRetryPolicy(new DefaultRetryPolicy(
                        20000,
                        DefaultRetryPolicy.DEFAULT_MAX_RETRIES,
                        DefaultRetryPolicy.DEFAULT_BACKOFF_MULT));
                queue.add(jsonObjectRequest);
            }
        });

    }
}
