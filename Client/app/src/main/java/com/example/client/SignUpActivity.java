package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.app.DatePickerDialog;
import android.content.Intent;
import android.graphics.Color;
import android.graphics.drawable.ColorDrawable;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.DatePicker;
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

import java.util.Calendar;

public class SignUpActivity extends AppCompatActivity {

    private TextView firstName;
    private TextView lastName;
    private TextView email;
    private TextView password;
    private TextView cardNumber;
    private TextView cardExpiryDate;
    private TextView cvv;
    private Button signUpButton;
    private DatePickerDialog.OnDateSetListener onDateSetListener;
    private String selectedDay;
    private String selectedMonth;
    private String selectedYear;
    private String creditCardExpiryDate;

    public void signUp(){
        boolean err = false;
        if(firstName.getText().length() == 0){
            err = true;
            firstName.setError("First Name is Required");
        }
        if(lastName.getText().length() == 0){
            err = true;
            lastName.setError("Last Name is Required");
        }
        if(!email.getText().toString().contains("@")){
            err = true;
            email.setError("Invalid Email");
        }
        if(password.getText().length() < 6){
            err = true;
            password.setError("Password should be 6 characters or more");
        }
        if(cardNumber.getText().length() < 6){
            err = true;
            cardNumber.setError("Invalid Card Number");
        }
        if(cardExpiryDate.getText().length() == 0){
            err = true;
            cardExpiryDate.setError("Card Expiry Date is required");
        }
        if(cvv.getText().length() == 0 || cvv.getText().length() > 4){
            err = true;
            cvv.setError("Invalid Card CVV");
        }
        if(!err) {
            RequestQueue queue = Volley.newRequestQueue(SignUpActivity.this);
            creditCardExpiryDate = selectedYear + "-" + selectedMonth + "-" + selectedDay;
            JSONObject jsonBody = new JSONObject();
            try {
                jsonBody.put("firstName", firstName.getText().toString());
                jsonBody.put("lastName", lastName.getText().toString());
                jsonBody.put("email", email.getText().toString());
                jsonBody.put("password", password.getText().toString());
                jsonBody.put("creditCardNumber", cardNumber.getText().toString());
                jsonBody.put("creditCardExpiryDate", creditCardExpiryDate);
                jsonBody.put("creditCardCVV", Integer.parseInt(cvv.getText().toString()));
            } catch (JSONException e) {
                e.printStackTrace();
            }
            String url = "http://10.0.2.2:8000/customers";
            JsonObjectRequest jsonObjectRequest = new JsonObjectRequest(Request.Method.POST, url, jsonBody, new Response.Listener<JSONObject>() {
                @Override
                public void onResponse(JSONObject response) {
                    Intent i = new Intent(getApplicationContext(), CustomerViewShopsActivity.class);
                    try {
                        i.putExtra("id", response.getInt("id"));
                        i.putExtra("firstName", response.getString("firstName"));
                        i.putExtra("lastName", response.getString("lastName"));
                        i.putExtra("email", response.getString("email"));
                        i.putExtra("password", response.getString("password"));
                        i.putExtra("creditCardNumber", response.getString("creditCardNumber"));
                        i.putExtra("creditCardExpiryDate", response.getString("creditCardExpiryDate"));
                        i.putExtra("creditCardCVV", response.getInt("creditCardCVV"));

                    } catch (JSONException e) {
                        e.printStackTrace();
                    }
                    startActivity(i);
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
    }
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_sign_up);

        firstName = findViewById(R.id.firstNameTextView);
        lastName = findViewById(R.id.lastNameTextView);
        email = findViewById(R.id.emailTextView);
        password = findViewById(R.id.passwordTextView);
        cardNumber = findViewById(R.id.cardNumberTextView);
        cardExpiryDate = findViewById(R.id.cardExpiryDateTextView);
        cvv = findViewById(R.id.cvvTextView);
        signUpButton = findViewById(R.id.customerSignUpButton);

        cardExpiryDate.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Calendar calendar = Calendar.getInstance();
                int year = calendar.get(Calendar.YEAR);
                int month = calendar.get(Calendar.MONTH);
                int day = calendar.get(Calendar.DAY_OF_MONTH);
                selectedDay = String.format("%02d", day);
                month++;
                selectedMonth = String.format("%02d", month);
                selectedYear = year+"";
                DatePickerDialog dialog = new DatePickerDialog(
                        SignUpActivity.this,
                        android.R.style.Theme_DeviceDefault_Dialog_MinWidth,
                        onDateSetListener,
                        year, month, day);
                dialog.getWindow().setBackgroundDrawable(new ColorDrawable(Color.TRANSPARENT));
                dialog.show();
            }
        });
        onDateSetListener = new DatePickerDialog.OnDateSetListener() {
            @Override
            public void onDateSet(DatePicker view, int year, int month, int dayOfMonth) {
                month++;
                selectedDay = String.format("%02d", dayOfMonth);
                selectedMonth = String.format("%02d", month);
                selectedYear = year+"";
                cardExpiryDate.setText(selectedYear+"/"+selectedMonth+"/"+selectedDay);
            }
        };
        signUpButton.setOnClickListener(new View.OnClickListener(){
            public void onClick(View v) {
                signUp();
            }
        });
    }
}
