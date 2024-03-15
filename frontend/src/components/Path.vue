<script setup>
import { defineProps, h, ref, onMounted, defineEmits, watch} from "vue";
import { useMessage, NIcon } from "naive-ui";
import { GetDirs } from '../../wailsjs/go/main/App';
import { Folder, FolderOpenOutline, FileTrayFullOutline , CloseOutline} from "@vicons/ionicons5";

const props = defineProps({
    placeholder: String,
    type: String,
})

const emits = defineEmits(["clickPath"]);
const message = useMessage();
const data = ref([])
const input = ref("")
const showClose = ref(false)

onMounted(async () => {
    try {
        const paths = await dirs("")
        if (paths.length != 0) {
            paths.forEach(path => {
                data.value.push({
                    key: path.path,
                    label: path.name,
                    prefix: () => h(NIcon, null, {
                        default: () => h(path.name ? Folder : FileTrayFullOutline)
                    }),
                    isDir: path.isDir,
                    disabled: props.type == "dir" && !path.isDir,
                    isLeaf: !path.hasChildren
                })
            })
        }
    } catch (e) {
        console.log(e)
    }
})

const clearInput = () => {
    input.value = ""
}

watch(input, (newValue, oldValue) => {
    newValue == "" ? showClose.value = false : showClose.value = true
})

const handleLoad = async (node) => {
    const paths = await dirs(node.key)
    return new Promise((resolve) => {
        node.children = [];
        paths.forEach(path => {
            node.children.push({
                key: path.path,
                label: path.name,
                prefix: () => h(NIcon, null, {
                    default: () => h(path.isDir ? Folder : FileTrayFullOutline)
                }),
                isDir: path.isDir,
                disabled: props.type == "dir" && !path.isDir,
                isLeaf: !path.hasChildren
            })
        })
        resolve();
    });
}

const dirs = async (path) => {
    return GetDirs(path).then(res => {
        if (res.code == 0) {
            return res.data
        } else {
            message.error(res.message)
            return []
        }
    }).catch(e => {
        console.log(e)
    })
}
const nodeInputProps = ({ option }) => {
    return {
        onClick() {
            if (option.isDir) {
                input.value = option.key

                // 通知父元素
                emits("clickPath", input.value)
            }
        },
        onContextmenu(e) {
            console.log(e.clientX, e.clientY);
            e.preventDefault();
        }
    };
};
const updatePrefixWithExpaned = (_keys, _option, meta) => {
    if (!meta.node)
        return;
    switch (meta.action) {
        case "expand":
            meta.node.prefix = () => h(NIcon, null, {
                default: () => h(FolderOpenOutline)
            });
            break;
        case "collapse":
            meta.node.prefix = () => h(NIcon, null, {
                default: () => h(Folder)
            });
            break;
    }
};
</script>

<template>
    <div class="select-path">
        <div class="close-input">
            <input class="custom-input" type="text" v-model="input" :placeholder="placeholder">
            <n-icon v-if="showClose" size="20" color="#0e7a0d" @click="clearInput">
                <close-outline />
            </n-icon>
        </div>
        <n-tree block-line expand-on-click virtual-scroll style="height: 220px" :on-load="handleLoad" :data="data"
            :node-props="nodeInputProps" :on-update:expanded-keys="updatePrefixWithExpaned" />
    </div>

</template>

<style scoped>
.select-path {
    padding: 0 10px 10px 10px;
}
.close-input{
    display: flex;
    align-items: center;
}
.close-input i{
    cursor: pointer;
    margin-bottom: 10px;
}
.custom-input {
    background: transparent;
    border: none;
    line-height: 30px;
    height: 30px;
    outline: none;
    border-bottom: 1px solid rgb(226, 218, 218);
    width: 100%;
    margin-bottom: 10px;
}
</style>