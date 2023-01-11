<template>
  <van-form @submit="onSubmit">
    <van-cell-group inset>
      <van-field
        v-model="pushdeer_url"
        name="pushdeer"
        label="pushdeer"
        disabled
      />
      <van-tag type="warning" size="large">一般不需要修改</van-tag>
      <van-grid>
        <van-grid-item>
          <van-cell center title="切换">
              <van-switch v-model="disabled" size="22px" />
          </van-cell>
        </van-grid-item >
      </van-grid>
      <van-field
        v-model="f.network"
        name="chain_id"
        label="chain id"
        :disabled="disabled"
      />
      <van-field
        v-model="f.rpc"
        name="rpc"
        label="rpc"
        :disabled="disabled"
      />

      <van-divider />
      <van-tag type="warning" size="success">需要修改的表单</van-tag>

      <van-field
        v-model="f.abi"
        name="abi"
        label="abi"
      />
      <van-field
        v-model="f.contract_address"
        name="contract_address"
        label="contract_address"
      />
      <van-field
        v-model="f.events"
        name="events"
        label="events"
      />
      <van-field
        placeholder="请输入 push key"
        v-model="f.push_key"
        name="push_key"
        label="push_key"
      />
    </van-cell-group>
    <div style="margin: 16px;">
      <van-button round block type="primary" native-type="submit" @click="submit">
        提交
      </van-button>
    </div>
  </van-form>

  <!-- <van-list
  finished-text="message list"
  style="margin-top: 20px"
  >
    <slot>
    message list
    </slot>
    <van-cell v-for="item in message_list" :key="item" title="" >
    </van-cell>
  </van-list> -->

</template>

<script>

import { ref } from 'vue';
const axios = require('axios')
import { showDialog } from 'vant';

export default {
  setup() {
    const disabled = ref(true);

    const network = '5';
    const pushdeer_url = 'http://43.206.196.16:8800'
    const abi = "[{\"anonymous\": false, \"inputs\": [{\"indexed\": true, \"internalType\": \"address\", \"name\": \"owner\", \"type\": \"address\"}, {\"indexed\": true, \"internalType\": \"address\", \"name\": \"spender\", \"type\": \"address\"}, {\"indexed\": false, \"internalType\": \"uint256\", \"name\": \"value\", \"type\": \"uint256\"}], \"name\": \"Approval\", \"type\": \"event\"}, {\"anonymous\": false, \"inputs\": [{\"indexed\": true, \"internalType\": \"address\", \"name\": \"from\", \"type\": \"address\"}, {\"indexed\": true, \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\"}, {\"indexed\": false, \"internalType\": \"uint256\", \"name\": \"value\", \"type\": \"uint256\"}], \"name\": \"Transfer\", \"type\": \"event\"}]";
    const rpc = "https://eth-goerli.g.alchemy.com/v2/5c8YbVzERLkTRMTYSotAYp07LR2cGbOk";
    const contract_address = "0x12bd33596139D016c7fB1b035Ed77668549f623c";
    const events = "Transfer(address,address,uint256)";
    const push_key = "";
    // const push_key = "PDU2TKCPwxU4FjhhH2Al78IQTYqna4ur491Oo";

    const message_list = ref(['a', 1, 2]);

    return {
      message_list: message_list,
      disabled: disabled,
      pushdeer_url,
      f: {
        network,
        abi,
        rpc,
        contract_address,
        events,
        push_key,
      }
    };
  },
  // created(){
  //   this.loadMessage()
  // },
  methods: {
    submit(){
      console.log(1)
      var request_form = JSON.parse(JSON.stringify(this.f))
      request_form.events = [request_form.events]
      axios.post(
          '/v1/task',
          {
              'method': 'collector_register',
              'params': [request_form]
          }
        ).then(r=>{
          if (!r.data.error) {
            showDialog({ message: '提交成功' });
            return
          }
          showDialog({ message: `参数错误, code: ${r.data.error.code} message: ${r.data.error.message}` });
        })
    },
    // loadMessage(){
    //   console.log(1111)
    //   axios.post(
    //     this.pushdeer_url + '/message/list',
    //   ).then(r=>{
    //     console.log(2222)
    //     console.log(r)
    //     setTimeout(
    //       this.loadMessage,
    //       500
    //     )
    //   })
    // }
  }
};

</script>