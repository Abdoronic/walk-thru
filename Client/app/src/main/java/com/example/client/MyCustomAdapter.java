package com.example.client;

import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.Button;
import android.widget.ListAdapter;
import android.widget.TextView;

import java.util.ArrayList;


public class MyCustomAdapter extends BaseAdapter implements ListAdapter {
    private ArrayList<String> list = new ArrayList<String>();
    private Context context;
    private int[]quantity;
    public MyCustomAdapter(ArrayList<String> list, Context context) {
        this.list = list;
        this.context = context;
        this.quantity = new int[list.size()];
        for(int i =0;i<this.quantity.length;i++){
            this.quantity[i]=0;
        }
    }
    public int[] getQuantity(){
        return quantity;
    }

    @Override
    public int getCount() {
        return list.size();
    }

    @Override
    public Object getItem(int pos) {
        return list.get(pos);
    }

    @Override
    public long getItemId(int pos) {
        return 0;
    }

    @Override
    public View getView(final int position, View convertView, ViewGroup parent) {
        View view = convertView;
        if (view == null) {
            LayoutInflater inflater = (LayoutInflater) context.getSystemService(Context.LAYOUT_INFLATER_SERVICE);
            view = inflater.inflate(R.layout.custom_items_list, null);
        }

        final TextView itemDescTextView = view.findViewById(R.id.itemDescTextView);
        itemDescTextView.setText(list.get(position));
        final TextView itemQuantityTextView = view.findViewById(R.id.itemQuantityTextView);
        Button removeButton = view.findViewById(R.id.removeButton);
        Button addButton = view.findViewById(R.id.addButton);

        removeButton.setOnClickListener(new View.OnClickListener(){
            @Override
            public void onClick(View v) {
                int oldQuantity = Integer.parseInt(itemQuantityTextView.getText().toString());
                oldQuantity--;
                if(oldQuantity>=0) {
                    quantity[position]=oldQuantity;
                    itemQuantityTextView.setText(oldQuantity+"");
                    notifyDataSetChanged();
                }
            }
        });
        addButton.setOnClickListener(new View.OnClickListener(){
            @Override
            public void onClick(View v) {
                itemQuantityTextView.setText((Integer.parseInt(itemQuantityTextView.getText().toString())+1)+"");
                quantity[position]=Integer.parseInt(itemQuantityTextView.getText().toString());
                notifyDataSetChanged();
            }
        });
        return view;
    }
}
