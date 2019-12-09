package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;
import android.widget.Toast;

import com.android.volley.Request;
import com.android.volley.RequestQueue;
import com.android.volley.Response;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.JsonObjectRequest;
import com.android.volley.toolbox.Volley;

import org.json.JSONException;
import org.json.JSONObject;

import java.util.HashMap;
import java.util.Map;

public class AdminActivity extends AppCompatActivity {
    private Button signInButton;
    private TextView usernameTextView;
    private TextView passwordTextView;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_admin);
        signInButton = findViewById(R.id.customerSignUpButton);
        usernameTextView = findViewById(R.id.emailTextView);
        passwordTextView = findViewById(R.id.passwordTextView);

        signInButton.setOnClickListener(new View.OnClickListener(){
            public void onClick(View v) {
                RequestQueue queue = Volley.newRequestQueue(AdminActivity.this);

                Map<String, String> params = new HashMap<String, String>();
                params.put("adminUsername", usernameTextView.getText().toString());
                params.put("adminPassword", passwordTextView.getText().toString());

                String url = "http://10.0.2.2:8000/shops/login";
                JsonObjectRequest jsonObjectRequest = new JsonObjectRequest(Request.Method.POST, url, new JSONObject(params), new Response.Listener<JSONObject>() {
                    @Override
                    public void onResponse(JSONObject response) {
                        Intent i = new Intent(getApplicationContext(), ShopViewItemsActivity.class);
                        try {
                            i.putExtra("id",response.getInt("id"));
                            i.putExtra("name",response.getString("name"));
                            i.putExtra("location",response.getString("location"));
                            i.putExtra("adminUsername",response.getString("adminUsername"));
                            i.putExtra("adminPassword",response.getString("adminPassword"));

                        } catch (JSONException e) {
                            e.printStackTrace();
                        }
                        startActivity(i);
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
                queue.add(jsonObjectRequest);
            }
        });

    }
}
