<script lang="ts">
    import { onMount } from 'svelte';

    // Task tipini tanımlayalım
    type Task = {
        id: number;
        name: string;
        completed: boolean;
    };

    let tasks: Task[] = [];
    let newTaskName: string = '';
    const backendUrl = 'http://localhost:8080/api'; // Backend API adresimiz

    // Sayfa yüklendiğinde görevleri backend'den çek
    onMount(async () => {
        const response = await fetch(`${backendUrl}/tasks`);
        if (response.ok) {
            tasks = await response.json();
        } else {
            console.error('Görevler yüklenemedi.');
        }
    });

    // Yeni görev ekleyen fonksiyon
    async function addTask() {
        if (!newTaskName.trim()) return;

        const response = await fetch(`${backendUrl}/tasks`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: newTaskName,
                completed: false,
            }),
        });

        if (response.ok) {
            const newTask = await response.json();
            tasks = [...tasks, newTask]; // Listeyi reaktif olarak güncelle
            newTaskName = ''; // Input'u temizle
        } else {
            console.error('Görev eklenemedi.');
        }
    }

    async function deleteTask(id: number) {
        const response = await fetch(`${backendUrl}/tasks/${id}`, {
            method: 'DELETE',
        });

        if (response.ok) {
            tasks = tasks.filter(task => task.id !== id); // Görevi listeden kaldır
        } else {
            console.error('Görev silinemedi.');
        }
    }

    async function updateTask(task: Task) {
        const response = await fetch(`${backendUrl}/tasks/${task.id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(task),
        });

        if (!response.ok) {
            console.error('Görev güncellenemedi.');
        }
    }

    async function getTask(id: number) {
        const response = await fetch(`${backendUrl}/tasks/${id}`);
        if (response.ok) {
            const task = await response.json();
            return task;
        } else {
            console.error('Görev yüklenemedi.');
            return null;
        }
    }
</script>


<main class="w-full max-w-md mx-auto p-4 bg-white shadow-md rounded-lg">
    <h1>Yapılacaklar Listesi</h1>

    <div class="mb-4">
        <input
            type="text"
            bind:value={newTaskName}
            placeholder="Yeni görev ekle..."
            class="border p-2 w-full"
        />
        <button on:click={addTask} class="bg-blue-500 text-white p-2 mt-2">Ekle</button>
    </div>
    <ul class="list-disc pl-5">
        {#each tasks as task (task.id)}
            <li class="flex items-center justify-between mb-2">
                <span class={task.completed ? 'line-through' : ''}>{task.name}</span>
                <div>
                    <button
                        on:click={() => {
                            task.completed = !task.completed;
                            updateTask(task);
                        }}
                        class="bg-green-500 text-white p-1 mr-2"
                    >
                        {task.completed ? 'Tamamlandı' : 'Tamamla'}
                    </button>
                    <button
                        on:click={() => deleteTask(task.id)}
                        class="bg-red-500 text-white p-1"
                    >
                        Sil
                    </button>
                </div>
            </li>
        {/each}
    </ul>
    {#if tasks.length === 0}
        <p class="text-gray-500">Henüz görev yok.</p>
    {/if}
</main>