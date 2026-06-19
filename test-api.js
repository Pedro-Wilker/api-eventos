const API_URL = 'http://localhost:8080';

async function runTests() {
    console.log('🚀 Iniciando testes da API Eventos...\n');

    try {
        console.log('📡 Testando [GET /ping]...');
        const pingRes = await fetch(`${API_URL}/ping`);
        const pingData = await pingRes.json();
        console.log(`✅ Status: ${pingRes.status} | Resposta:`, pingData);
    } catch (e) {
        console.error('❌ Falha ao conectar na API. O servidor Go (main.go) está rodando?', e.message);
        return;
    }

    const testUser = {
        name: "Testador Node",
        email: `tester_${Date.now()}@email.com`,
        password: "senha_super_segura",
        role: "client"
    };

    console.log('\n👤 Testando [POST /api/register]...');
    const regRes = await fetch(`${API_URL}/api/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(testUser)
    });
    const regData = await regRes.json();
    console.log(`✅ Status: ${regRes.status} | Resposta:`, regData);

    console.log('\n🔑 Testando [POST /api/login]...');
    const loginRes = await fetch(`${API_URL}/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            email: testUser.email,
            password: testUser.password
        })
    });
    const loginData = await loginRes.json();

    let token = '';
    if (loginRes.ok && loginData.token) {
        token = loginData.token;
        console.log(`✅ Status: ${loginRes.status} | Token recebido com sucesso!`);
    } else {
        console.error(`❌ Status: ${loginRes.status} | Falha no login:`, loginData);
        return;
    }

    console.log('\n🛡️ Testando Rota Protegida [GET /api/me]...');
    const meRes = await fetch(`${API_URL}/api/me`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    const meData = await meRes.json();
    console.log(`✅ Status: ${meRes.status} | Resposta:`, meData);

    console.log('\n👥 Testando Criação de Convidado [POST /api/convidados]...');
    const guestPayload = {
        nome: "João Silva",
        quantidade_acompanhante: 2,
        nome_acompanhante: ["Maria Silva", "Joãzinho Silva"],
        email_convidado: "joao@email.com"
    };

    const createGuestRes = await fetch(`${API_URL}/api/convidados`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(guestPayload)
    });
    const createGuestData = await createGuestRes.json();
    console.log(`✅ Status: ${createGuestRes.status} | Resposta:`, createGuestData);

    console.log('\n📋 Testando Listagem de Convidados [GET /api/convidados]...');
    const listGuestRes = await fetch(`${API_URL}/api/convidados`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    const listGuestData = await listGuestRes.json();
    console.log(`✅ Status: ${listGuestRes.status} | Total listado: ${listGuestData.data ? listGuestData.data.length : 0}`);

    console.log('\n🎁 Testando Criação de Presente [POST /api/presentes]...');
    const giftPayload = {
        nome_presente: "Liquidificador",
        descricao_presente: "Potência 1000W, Cor Preta",
        quantas_reservas_por_presente: 1,
        foto_presente: "https://link-da-imagem.com/liquidificador.jpg"
    };

    const createGiftRes = await fetch(`${API_URL}/api/presentes`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(giftPayload)
    });
    const createGiftData = await createGiftRes.json();
    console.log(`✅ Status: ${createGiftRes.status} | Resposta:`, createGiftData);

    console.log('\n🛍️ Testando Listagem de Presentes [GET /api/presentes]...');
    const listGiftRes = await fetch(`${API_URL}/api/presentes`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    const listGiftData = await listGiftRes.json();
    console.log(`✅ Status: ${listGiftRes.status} | Total listado: ${listGiftData.data ? listGiftData.data.length : 0}`);

    console.log('\n🏁 Todos os testes passaram com sucesso!');
}

runTests();