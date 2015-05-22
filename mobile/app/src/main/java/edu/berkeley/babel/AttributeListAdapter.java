package edu.berkeley.babel;

import android.content.Context;
import android.os.Looper;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.EditText;
import android.widget.TextView;

import java.util.LinkedList;
import java.util.List;

public class AttributeListAdapter extends BaseAdapter {
    private Context mContext;
    private boolean mEnabled;
    private List<Pair<String, String>> mAttributes;

    public static class Pair<F, S> { // self-defined mutable pair
        public F first;
        public S second;

        public Pair(F first, S second) {
            this.first = first;
            this.second = second;
        }
    }

    private static class ViewHolder {
        public TextView mNameView;
        public EditText mValueView;
    }

    class ValueFocusChangeListener implements View.OnFocusChangeListener {
        private int mPosition;

        public ValueFocusChangeListener(int position) {
            mPosition = position;
        }


        @Override
        public void onFocusChange(View v, boolean hasFocus) {
            if (hasFocus) {
                return;
            }
            // update if not focused
            Pair<String, String> attribute = mAttributes.get(mPosition);
            EditText valueView = (EditText) v;
            attribute.second = valueView.getText().toString();
        }
    }

    public AttributeListAdapter(Context context) {
        mContext = context;
        mEnabled = true;
        mAttributes = new LinkedList<>();
    }

    public void addAll(List<Pair<String, String>> attributes) {
        // only allow updates from main thread
        // ref: http://www.piwai.info/android-adapter-good-practices/#Thread-safety
        if (BuildConfig.DEBUG) {
            if (Thread.currentThread() != Looper.getMainLooper().getThread()) {
                throw new IllegalStateException("This method should be called from the Main Thread");
            }
        }

        mAttributes = attributes;
        notifyDataSetChanged();
    }

    public void add(Pair<String, String> attribute) {
        // only allow updates from main thread
        // ref: http://www.piwai.info/android-adapter-good-practices/#Thread-safety
        if (BuildConfig.DEBUG) {
            if (Thread.currentThread() != Looper.getMainLooper().getThread()) {
                throw new IllegalStateException("This method should be called from the Main Thread");
            }
        }

        mAttributes.add(attribute);
        notifyDataSetChanged();
    }

    public void clear() {
        mAttributes.clear();
    }

    public void setEnabled(boolean enabled) {
        mEnabled = enabled;
        super.notifyDataSetChanged();
    }

    @Override
    public int getCount() {
        return mAttributes.size();
    }

    @Override
    public Pair<String, String> getItem(int position) {
        return mAttributes.get(position);
    }

    @Override
    public long getItemId(int position) {
        return position;
    }

    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        View itemView;

        if (convertView == null) {
            LayoutInflater inflater = (LayoutInflater) mContext.getSystemService(Context.LAYOUT_INFLATER_SERVICE);
            itemView = inflater.inflate(R.layout.attribute_item, null);

            ViewHolder newHolder = new ViewHolder();
            newHolder.mNameView = (TextView) itemView.findViewById(R.id.name);
            newHolder.mValueView = (EditText) itemView.findViewById(R.id.value);
            itemView.setTag(newHolder);
        } else {
            itemView = convertView;
        }

        ViewHolder holder = (ViewHolder) itemView.getTag();
        Pair<String, String> attr = mAttributes.get(position);
        holder.mNameView.setText(attr.first);
        holder.mNameView.setEnabled(mEnabled);
        holder.mValueView.setText(attr.second);
        holder.mValueView.setEnabled(mEnabled);
        holder.mValueView.setOnFocusChangeListener(new ValueFocusChangeListener(position));

        return itemView;
    }
}
