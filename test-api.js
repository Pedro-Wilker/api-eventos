const API_URL = 'https://api-eventos.artonbyte.com.br';

// Pequeno helper pra deixar o log mais legível
function logStep(emoji, title) {
    console.log(`\n${emoji} ${title}`);
}

function logResult(status, data) {
    console.log(`   Status: ${status} | Resposta:`, data);
}

async function runTests() {
    console.log('🚀 Iniciando testes da API Eventos...\n');

    // ----------------------------------------------------------------------
    // 0. Health Check
    // ----------------------------------------------------------------------
    try {
        logStep('📡', 'Testando [GET /ping]...');
        const pingRes = await fetch(`${API_URL}/ping`);
        const pingData = await pingRes.json();
        logResult(pingRes.status, pingData);
    } catch (e) {
        console.error('❌ Falha ao conectar na API. O servidor Go (main.go) está rodando?', e.message);
        return;
    }

    // ----------------------------------------------------------------------
    // 1. Registro e Login (client)
    // ----------------------------------------------------------------------
    const testUser = {
        name: "Testador Node",
        email: `tester_${Date.now()}@email.com`,
        password: "senha_super_segura",
        role: "client"
    };

    logStep('👤', 'Testando [POST /api/register]...');
    const regRes = await fetch(`${API_URL}/api/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(testUser)
    });
    const regData = await regRes.json();
    logResult(regRes.status, regData);

    logStep('🔑', 'Testando [POST /api/login]...');
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
        console.log(`   Status: ${loginRes.status} | Token recebido com sucesso!`);
    } else {
        console.error(`❌ Status: ${loginRes.status} | Falha no login:`, loginData);
        return;
    }

    logStep('🔒', 'Testando [POST /api/login] com senha errada (esperado 401)...');
    const badLoginRes = await fetch(`${API_URL}/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            email: testUser.email,
            password: "senha_errada"
        })
    });
    const badLoginData = await badLoginRes.json();
    logResult(badLoginRes.status, badLoginData);

    // ----------------------------------------------------------------------
    // 2. Perfil do usuário autenticado
    // ----------------------------------------------------------------------
    logStep('🛡️', 'Testando Rota Protegida [GET /api/me]...');
    const meRes = await fetch(`${API_URL}/api/me`, {
        method: 'GET',
        headers: { 'Authorization': `Bearer ${token}` }
    });
    const meData = await meRes.json();
    logResult(meRes.status, meData);
    const userId = meData.user_id;

    // ----------------------------------------------------------------------
    // 3. Convidados — CRUD
    // ----------------------------------------------------------------------
    logStep('👥', 'Testando Criação de Convidado [POST /api/convidados]...');
    const guestPayload = {
        nome: "João Silva",
        quantidade_acompanhante: 2,
        nome_acompanhante: ["Maria Silva", "Joãzinho Silva"],
        email_convidado: "joao@email.com",
        numero_convidado: "11999999999",
        emails_acompanhantes: ["maria@email.com"],
        numeros_acompanhantes: ["11888888888"]
        // ⚠️ "relacoes_acompanhante" removido temporariamente: a API retornou
        // 500 "Falha ao salvar convidado" com esse campo em produção, o que
        // sugere que a coluna ainda não foi migrada no banco. Reative quando
        // confirmar que a migração foi aplicada.
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
    logResult(createGuestRes.status, createGuestData);

    const guest = createGuestData.data || {};
    const guestId = guest.id ?? guest.ID;
    const guestQrCode = guest.qr_code ?? guest.QrCode ?? guest.QRCode;

    logStep('📋', 'Testando Listagem de Convidados [GET /api/convidados]...');
    const listGuestRes = await fetch(`${API_URL}/api/convidados`, {
        method: 'GET',
        headers: { 'Authorization': `Bearer ${token}` }
    });
    const listGuestData = await listGuestRes.json();
    console.log(`   Status: ${listGuestRes.status} | Total listado: ${listGuestData.data ? listGuestData.data.length : 0}`);

    logStep('✏️', `Testando Atualização de Convidado [PUT /api/convidados/${guestId}]...`);
    // ⚠️ Assim como em Presentes, esse PUT substitui o objeto inteiro — campos
    // omitidos (telefone, e-mail, acompanhantes...) seriam apagados. Aqui
    // testamos só com os campos básicos de propósito, pra documentar esse
    // comportamento; em uso real, envie o objeto completo.
    const updateGuestRes = await fetch(`${API_URL}/api/convidados/${guestId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
            nome: "João Silva Atualizado",
            quantidade_acompanhante: 3
        })
    });
    const updateGuestData = await updateGuestRes.json();
    logResult(updateGuestRes.status, updateGuestData);

    logStep('🔍', `Testando Buscar Convidado por Código (QR Code) [GET /api/convidados/buscar/${guestQrCode}]...`);
    const buscaRes = await fetch(`${API_URL}/api/convidados/buscar/${guestQrCode}`, {
        method: 'GET',
        headers: { 'Authorization': `Bearer ${token}` }
    });
    const buscaData = await buscaRes.json();
    logResult(buscaRes.status, buscaData);

    logStep('✅', 'Testando Check-in do Convidado [POST /api/convidados/checkin]...');
    const checkinRes = await fetch(`${API_URL}/api/convidados/checkin`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ codigo: guestQrCode })
    });
    const checkinData = await checkinRes.json();
    logResult(checkinRes.status, checkinData);

    logStep('♻️', 'Testando Check-in Duplicado (esperado 409)...');
    const checkinDupRes = await fetch(`${API_URL}/api/convidados/checkin`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ codigo: guestQrCode })
    });
    const checkinDupData = await checkinDupRes.json();
    logResult(checkinDupRes.status, checkinDupData);

    logStep('🚫', 'Testando Check-in com código inexistente (esperado 404)...');
    const checkinInvalidoRes = await fetch(`${API_URL}/api/convidados/checkin`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ codigo: "00000000-0000-0000-0000-000000000000" })
    });
    const checkinInvalidoData = await checkinInvalidoRes.json();
    logResult(checkinInvalidoRes.status, checkinInvalidoData);

    // ----------------------------------------------------------------------
    // 4. Presentes — CRUD
    // ----------------------------------------------------------------------
    logStep('🎁', 'Testando Criação de Presente [POST /api/presentes]...');
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
    logResult(createGiftRes.status, createGiftData);
    const giftRecord = createGiftData.data || createGiftData;
    const giftId = giftRecord.id ?? giftRecord.ID;

    logStep('🛍️', 'Testando Listagem de Presentes [GET /api/presentes]...');
    const listGiftRes = await fetch(`${API_URL}/api/presentes`, {
        method: 'GET',
        headers: { 'Authorization': `Bearer ${token}` }
    });
    const listGiftData = await listGiftRes.json();
    console.log(`   Status: ${listGiftRes.status} | Total listado: ${listGiftData.data ? listGiftData.data.length : (Array.isArray(listGiftData) ? listGiftData.length : 0)}`);

    logStep('✏️', `Testando Atualização de Presente [PUT /api/presentes/${giftId}]...`);
    // ⚠️ O PUT substitui o objeto inteiro (não faz merge parcial) — por isso
    // enviamos TODOS os campos aqui, mesmo os que não mudaram, pra não apagar
    // descrição/foto. Veja "Observações de Implementação" no README.
    // Mantemos a capacidade em 1 propositalmente, pra testar o esgotamento
    // (409) logo abaixo com uma única reserva.
    const updateGiftRes = await fetch(`${API_URL}/api/presentes/${giftId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
            nome_presente: "Liquidificador Atualizado",
            descricao_presente: "Potência 1000W, Cor Preta",
            quantas_reservas_por_presente: 1,
            foto_presente: "https://link-da-imagem.com/liquidificador.jpg"
        })
    });
    const updateGiftData = await updateGiftRes.json();
    logResult(updateGiftRes.status, updateGiftData);

    // ----------------------------------------------------------------------
    // 5. Rotas públicas de presentes (página do convite)
    // ----------------------------------------------------------------------
    logStep('🌐', `Testando Listagem Pública de Presentes [GET /api/public/presentes/${userId}]...`);
    const publicGiftsRes = await fetch(`${API_URL}/api/public/presentes/${userId}`);
    const publicGiftsData = await publicGiftsRes.json();
    logResult(publicGiftsRes.status, publicGiftsData);

    logStep('🎉', `Testando Reserva de Presente [POST /api/public/presentes/${giftId}/reservar]...`);
    // Capacidade foi fixada em 1 no PUT acima, então esta primeira reserva
    // deve ter sucesso (200) e a próxima deve esgotar (409).
    const reservarRes = await fetch(`${API_URL}/api/public/presentes/${giftId}/reservar`, {
        method: 'POST'
    });
    const reservarData = await reservarRes.json();
    logResult(reservarRes.status, reservarData);

    logStep('🚫', 'Testando Reserva de Presente já esgotado (esperado 409)...');
    const reservarEsgotadoRes = await fetch(`${API_URL}/api/public/presentes/${giftId}/reservar`, {
        method: 'POST'
    });
    const reservarEsgotadoData = await reservarEsgotadoRes.json();
    logResult(reservarEsgotadoRes.status, reservarEsgotadoData);

    // ----------------------------------------------------------------------
    // 6. Limpeza — deletar convidado e presente criados no teste
    // ----------------------------------------------------------------------
    logStep('🗑️', `Testando Deleção de Convidado [DELETE /api/convidados/${guestId}]...`);
    const deleteGuestRes = await fetch(`${API_URL}/api/convidados/${guestId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
    });
    const deleteGuestData = await deleteGuestRes.json();
    logResult(deleteGuestRes.status, deleteGuestData);

    logStep('🗑️', `Testando Deleção de Presente [DELETE /api/presentes/${giftId}]...`);
    const deleteGiftRes = await fetch(`${API_URL}/api/presentes/${giftId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
    });
    const deleteGiftData = await deleteGiftRes.json();
    logResult(deleteGiftRes.status, deleteGiftData);

    console.log('\n🏁 Todos os testes foram executados!');
}

runTests();