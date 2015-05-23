package edu.berkeley.babel.util;

import android.content.Context;
import android.os.Looper;
import android.text.Editable;
import android.text.TextWatcher;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.EditText;
import android.widget.TextView;

import java.util.LinkedList;
import java.util.List;

import edu.berkeley.babel.BuildConfig;
import edu.berkeley.babel.R;

public class KeyValueListAdapter extends BaseAdapter {
    private Context mContext;
    private OnKeyValueChangedListener mListener;
    private boolean mEnabled;
    private List<Pair<String, String>> mKeyValues;

    public static class Pair<F, S> { // self-defined mutable pair
        public F first;
        public S second;

        public Pair(F first, S second) {
            this.first = first;
            this.second = second;
        }
    }

    public interface OnKeyValueChangedListener {
        public abstract void OnValueChanged(String key, String newValue);
    }

    private static class ViewHolder {
        public TextView mKeyView;
        public EditText mValueView;
    }

    class ValueTextWatcher implements TextWatcher {
        private int mPosition;

        public ValueTextWatcher(int position) {
            mPosition = position;
        }

        @Override
        public void onTextChanged(CharSequence s, int start, int before, int count) {
        }

        @Override
        public void beforeTextChanged(CharSequence s, int start, int count, int after) {
        }

        @Override
        public void afterTextChanged(Editable s) {
            Pair<String, String> keyValue = mKeyValues.get(mPosition);
            keyValue.second = s.toString();

            if (mListener != null) {
                mListener.OnValueChanged(keyValue.first, s.toString());
            }

        }
    }

    public KeyValueListAdapter(Context context) {
        mContext = context;
        mEnabled = true;
        mKeyValues = new LinkedList<>();
    }

    public void setOnKeyValueChangedListener(OnKeyValueChangedListener listener) {
        mListener = listener;
    }

    public void add(Pair<String, String> keyValue) {
        // only allow updates from main thread
        // ref: http://www.piwai.info/android-adapter-good-practices/#Thread-safety
        if (BuildConfig.DEBUG) {
            if (Thread.currentThread() != Looper.getMainLooper().getThread()) {
                throw new IllegalStateException("This method should be called from the Main Thread");
            }
        }

        mKeyValues.add(keyValue);
    }

    public void clear() {
        mKeyValues.clear();
    }

    public void setEnabled(boolean enabled) {
        mEnabled = enabled;
        notifyDataSetChanged();
    }

    @Override
    public int getCount() {
        return mKeyValues.size();
    }

    @Override
    public Pair<String, String> getItem(int position) {
        return mKeyValues.get(position);
    }

    @Override
    public long getItemId(int position) {
        return position;
    }

    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        // not using convertView because i don't want to deal with removing TextChangedListener of old valueView
        mListener = null;

        LayoutInflater inflater = (LayoutInflater) mContext.getSystemService(Context.LAYOUT_INFLATER_SERVICE);
        View itemView = inflater.inflate(R.layout.key_value_item, null);

        ViewHolder newHolder = new ViewHolder();
        newHolder.mKeyView = (TextView) itemView.findViewById(R.id.key);
        newHolder.mValueView = (EditText) itemView.findViewById(R.id.value);
        itemView.setTag(newHolder);


        ViewHolder holder = (ViewHolder) itemView.getTag();
        Pair<String, String> keyValue = mKeyValues.get(position);
        holder.mKeyView.setText(keyValue.first);
        holder.mKeyView.setEnabled(mEnabled);
        holder.mValueView.setText(keyValue.second);
        holder.mValueView.setEnabled(mEnabled);
        holder.mValueView.addTextChangedListener(new ValueTextWatcher(position));

        return itemView;
    }
}
